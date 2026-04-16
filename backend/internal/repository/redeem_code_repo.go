package repository

import (
	"context"
	"sort"
	"strings"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/redeemcode"
	"github.com/Wei-Shaw/sub2api/ent/redeemcodeclaim"
	"github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"

	entsql "entgo.io/ent/dialect/sql"
)

type redeemCodeRepository struct {
	client *dbent.Client
}

func NewRedeemCodeRepository(client *dbent.Client) service.RedeemCodeRepository {
	return &redeemCodeRepository{client: client}
}

func (r *redeemCodeRepository) Create(ctx context.Context, code *service.RedeemCode) error {
	client := clientFromContext(ctx, r.client)
	builder := client.RedeemCode.Create().
		SetCode(code.Code).
		SetType(code.Type).
		SetValue(code.Value).
		SetStatus(code.Status).
		SetNotes(code.Notes).
		SetValidityDays(code.ValidityDays).
		SetNillableUsedBy(code.UsedBy).
		SetNillableUsedAt(code.UsedAt).
		SetNillableGroupID(code.GroupID)
	if code.MaxClaims > 0 {
		builder.SetMaxClaims(code.MaxClaims)
	}
	if code.ClaimedCount > 0 {
		builder.SetClaimedCount(code.ClaimedCount)
	}
	created, err := builder.Save(ctx)
	err = translatePersistenceError(err, nil, service.ErrRedeemCodeExists)
	if err == nil {
		code.ID = created.ID
		code.CreatedAt = created.CreatedAt
		code.MaxClaims = created.MaxClaims
		code.ClaimedCount = created.ClaimedCount
	}
	return err
}

func (r *redeemCodeRepository) CreateBatch(ctx context.Context, codes []service.RedeemCode) error {
	if len(codes) == 0 {
		return nil
	}

	builders := make([]*dbent.RedeemCodeCreate, 0, len(codes))
	client := clientFromContext(ctx, r.client)
	for i := range codes {
		c := &codes[i]
		b := client.RedeemCode.Create().
			SetCode(c.Code).
			SetType(c.Type).
			SetValue(c.Value).
			SetStatus(c.Status).
			SetNotes(c.Notes).
			SetValidityDays(c.ValidityDays).
			SetNillableUsedBy(c.UsedBy).
			SetNillableUsedAt(c.UsedAt).
			SetNillableGroupID(c.GroupID)
		if c.MaxClaims > 0 {
			b.SetMaxClaims(c.MaxClaims)
		}
		if c.ClaimedCount > 0 {
			b.SetClaimedCount(c.ClaimedCount)
		}
		builders = append(builders, b)
	}

	err := client.RedeemCode.CreateBulk(builders...).Exec(ctx)
	return translatePersistenceError(err, nil, service.ErrRedeemCodeExists)
}

func (r *redeemCodeRepository) GetByID(ctx context.Context, id int64) (*service.RedeemCode, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.RedeemCode.Query().
		Where(redeemcode.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrRedeemCodeNotFound
		}
		return nil, err
	}
	return redeemCodeEntityToService(m), nil
}

func (r *redeemCodeRepository) GetByCode(ctx context.Context, code string) (*service.RedeemCode, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.RedeemCode.Query().
		Where(redeemcode.CodeEQ(code)).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrRedeemCodeNotFound
		}
		return nil, err
	}
	return redeemCodeEntityToService(m), nil
}

func (r *redeemCodeRepository) Delete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.RedeemCode.Delete().Where(redeemcode.IDEQ(id)).Exec(ctx)
	return err
}

func (r *redeemCodeRepository) List(ctx context.Context, params pagination.PaginationParams) ([]service.RedeemCode, *pagination.PaginationResult, error) {
	return r.ListWithFilters(ctx, params, "", "", "")
}

func (r *redeemCodeRepository) ListWithFilters(ctx context.Context, params pagination.PaginationParams, codeType, status, search string) ([]service.RedeemCode, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := client.RedeemCode.Query()

	if codeType != "" {
		q = q.Where(redeemcode.TypeEQ(codeType))
	}
	if status != "" {
		q = q.Where(redeemcode.StatusEQ(status))
	}
	if search != "" {
		q = q.Where(
			redeemcode.Or(
				redeemcode.CodeContainsFold(search),
				redeemcode.HasUserWith(user.EmailContainsFold(search)),
			),
		)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	codesQuery := q.
		WithUser().
		WithGroup().
		Offset(params.Offset()).
		Limit(params.Limit())
	for _, order := range redeemCodeListOrder(params) {
		codesQuery = codesQuery.Order(order)
	}

	codes, err := codesQuery.All(ctx)
	if err != nil {
		return nil, nil, err
	}

	outCodes := redeemCodeEntitiesToService(codes)

	return outCodes, paginationResultFromTotal(int64(total), params), nil
}

func redeemCodeListOrder(params pagination.PaginationParams) []func(*entsql.Selector) {
	sortBy := strings.ToLower(strings.TrimSpace(params.SortBy))
	sortOrder := params.NormalizedSortOrder(pagination.SortOrderDesc)

	var field string
	switch sortBy {
	case "type":
		field = redeemcode.FieldType
	case "value":
		field = redeemcode.FieldValue
	case "status":
		field = redeemcode.FieldStatus
	case "used_at":
		field = redeemcode.FieldUsedAt
	case "created_at":
		field = redeemcode.FieldCreatedAt
	case "code":
		field = redeemcode.FieldCode
	default:
		field = redeemcode.FieldID
	}

	if sortOrder == pagination.SortOrderAsc {
		return []func(*entsql.Selector){dbent.Asc(field), dbent.Asc(redeemcode.FieldID)}
	}
	return []func(*entsql.Selector){dbent.Desc(field), dbent.Desc(redeemcode.FieldID)}
}

func (r *redeemCodeRepository) Update(ctx context.Context, code *service.RedeemCode) error {
	client := clientFromContext(ctx, r.client)
	up := client.RedeemCode.UpdateOneID(code.ID).
		SetCode(code.Code).
		SetType(code.Type).
		SetValue(code.Value).
		SetStatus(code.Status).
		SetNotes(code.Notes).
		SetValidityDays(code.ValidityDays).
		SetMaxClaims(code.MaxClaims).
		SetClaimedCount(code.ClaimedCount)

	if code.UsedBy != nil {
		up.SetUsedBy(*code.UsedBy)
	} else {
		up.ClearUsedBy()
	}
	if code.UsedAt != nil {
		up.SetUsedAt(*code.UsedAt)
	} else {
		up.ClearUsedAt()
	}
	if code.GroupID != nil {
		up.SetGroupID(*code.GroupID)
	} else {
		up.ClearGroupID()
	}

	updated, err := up.Save(ctx)
	err = translatePersistenceError(err, service.ErrRedeemCodeNotFound, service.ErrRedeemCodeExists)
	if err != nil {
		return err
	}
	code.CreatedAt = updated.CreatedAt
	return nil
}

func (r *redeemCodeRepository) Use(ctx context.Context, id, userID int64) error {
	now := time.Now()
	client := clientFromContext(ctx, r.client)
	affected, err := client.RedeemCode.Update().
		Where(redeemcode.IDEQ(id), redeemcode.StatusEQ(service.StatusUnused)).
		SetStatus(service.StatusUsed).
		SetUsedBy(userID).
		SetUsedAt(now).
		Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return service.ErrRedeemCodeUsed
	}
	return nil
}

func (r *redeemCodeRepository) HasClaimByUser(ctx context.Context, redeemCodeID, userID int64) (bool, error) {
	client := clientFromContext(ctx, r.client)
	return client.RedeemCodeClaim.Query().
		Where(
			redeemcodeclaim.RedeemCodeIDEQ(redeemCodeID),
			redeemcodeclaim.UserIDEQ(userID),
		).
		Exist(ctx)
}

func (r *redeemCodeRepository) CreateClaim(ctx context.Context, redeemCodeID, userID int64, amount float64) error {
	client := clientFromContext(ctx, r.client)
	err := client.RedeemCodeClaim.Create().
		SetRedeemCodeID(redeemCodeID).
		SetUserID(userID).
		SetAmount(amount).
		Exec(ctx)
	return translatePersistenceError(err, nil, service.ErrRedeemCodeClaimed)
}

func (r *redeemCodeRepository) IncrementClaimedCount(ctx context.Context, id, maxClaims int64) (int, error) {
	client := clientFromContext(ctx, r.client)
	affected, err := client.RedeemCode.Update().
		Where(
			redeemcode.IDEQ(id),
			redeemcode.StatusEQ(service.StatusUnused),
			redeemcode.ClaimedCountLT(int(maxClaims)),
		).
		AddClaimedCount(1).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	if affected == 0 {
		return 0, service.ErrRedeemCodeUsed
	}

	updated, err := client.RedeemCode.Query().
		Where(redeemcode.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return 0, service.ErrRedeemCodeNotFound
		}
		return 0, err
	}
	return updated.ClaimedCount, nil
}

func (r *redeemCodeRepository) ListByUser(ctx context.Context, userID int64, limit int) ([]service.RedeemCode, error) {
	if limit <= 0 {
		limit = 10
	}

	client := clientFromContext(ctx, r.client)
	codes, err := client.RedeemCode.Query().
		Where(redeemcode.UsedByEQ(userID)).
		WithGroup().
		Order(dbent.Desc(redeemcode.FieldUsedAt)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := client.RedeemCodeClaim.Query().
		Where(redeemcodeclaim.UserIDEQ(userID)).
		WithRedeemCode(func(q *dbent.RedeemCodeQuery) {
			q.WithGroup()
		}).
		Order(dbent.Desc(redeemcodeclaim.FieldClaimedAt)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}

	combined := append(redeemCodeEntitiesToService(codes), redeemCodeClaimEntitiesToService(claims)...)
	sortRedeemCodesByUsedAtDesc(combined)
	if len(combined) > limit {
		combined = combined[:limit]
	}
	return combined, nil
}

// ListByUserPaginated returns paginated balance/concurrency history for a user.
// Supports optional type filter (e.g. "balance", "admin_balance", "concurrency", "admin_concurrency", "subscription").
func (r *redeemCodeRepository) ListByUserPaginated(ctx context.Context, userID int64, params pagination.PaginationParams, codeType string) ([]service.RedeemCode, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := client.RedeemCode.Query().
		Where(redeemcode.UsedByEQ(userID))

	// Optional type filter
	if codeType != "" {
		q = q.Where(redeemcode.TypeEQ(codeType))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	codes, err := q.
		WithGroup().
		Order(dbent.Desc(redeemcode.FieldUsedAt)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	claimQuery := client.RedeemCodeClaim.Query().
		Where(redeemcodeclaim.UserIDEQ(userID))
	if codeType != "" {
		claimQuery = claimQuery.Where(redeemcodeclaim.HasRedeemCodeWith(redeemcode.TypeEQ(codeType)))
	}
	claimTotal, err := claimQuery.Count(ctx)
	if err != nil {
		return nil, nil, err
	}
	claims, err := claimQuery.
		WithRedeemCode(func(q *dbent.RedeemCodeQuery) {
			q.WithGroup()
		}).
		Order(dbent.Desc(redeemcodeclaim.FieldClaimedAt)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	combined := append(redeemCodeEntitiesToService(codes), redeemCodeClaimEntitiesToService(claims)...)
	sortRedeemCodesByUsedAtDesc(combined)
	start := params.Offset()
	if start > len(combined) {
		start = len(combined)
	}
	end := start + params.Limit()
	if end > len(combined) {
		end = len(combined)
	}

	return combined[start:end], paginationResultFromTotal(int64(total+claimTotal), params), nil
}

// SumPositiveBalanceByUser returns total recharged amount (sum of value > 0 where type is balance/admin_balance).
func (r *redeemCodeRepository) SumPositiveBalanceByUser(ctx context.Context, userID int64) (float64, error) {
	client := clientFromContext(ctx, r.client)
	var result []struct {
		Sum float64 `json:"sum"`
	}
	err := client.RedeemCode.Query().
		Where(
			redeemcode.UsedByEQ(userID),
			redeemcode.ValueGT(0),
			redeemcode.TypeIn("balance", "admin_balance"),
		).
		Aggregate(dbent.As(dbent.Sum(redeemcode.FieldValue), "sum")).
		Scan(ctx, &result)
	if err != nil {
		return 0, err
	}

	var welfareResult []struct {
		Sum float64 `json:"sum"`
	}
	err = client.RedeemCodeClaim.Query().
		Where(
			redeemcodeclaim.UserIDEQ(userID),
			redeemcodeclaim.AmountGT(0),
			redeemcodeclaim.HasRedeemCodeWith(redeemcode.TypeEQ(service.RedeemTypeWelfare)),
		).
		Aggregate(dbent.As(dbent.Sum(redeemcodeclaim.FieldAmount), "sum")).
		Scan(ctx, &welfareResult)
	if err != nil {
		return 0, err
	}

	var total float64
	if len(result) > 0 {
		total += result[0].Sum
	}
	if len(welfareResult) > 0 {
		total += welfareResult[0].Sum
	}
	return total, nil
}

func redeemCodeEntityToService(m *dbent.RedeemCode) *service.RedeemCode {
	if m == nil {
		return nil
	}
	out := &service.RedeemCode{
		ID:           m.ID,
		Code:         m.Code,
		Type:         m.Type,
		Value:        m.Value,
		Status:       m.Status,
		UsedBy:       m.UsedBy,
		UsedAt:       m.UsedAt,
		Notes:        derefString(m.Notes),
		CreatedAt:    m.CreatedAt,
		MaxClaims:    m.MaxClaims,
		ClaimedCount: m.ClaimedCount,
		GroupID:      m.GroupID,
		ValidityDays: m.ValidityDays,
	}
	if m.Edges.User != nil {
		out.User = userEntityToService(m.Edges.User)
	}
	if m.Edges.Group != nil {
		out.Group = groupEntityToService(m.Edges.Group)
	}
	return out
}

func redeemCodeEntitiesToService(models []*dbent.RedeemCode) []service.RedeemCode {
	out := make([]service.RedeemCode, 0, len(models))
	for i := range models {
		if s := redeemCodeEntityToService(models[i]); s != nil {
			out = append(out, *s)
		}
	}
	return out
}

func redeemCodeClaimEntitiesToService(models []*dbent.RedeemCodeClaim) []service.RedeemCode {
	out := make([]service.RedeemCode, 0, len(models))
	for i := range models {
		model := models[i]
		if model == nil || model.Edges.RedeemCode == nil {
			continue
		}

		parent := redeemCodeEntityToService(model.Edges.RedeemCode)
		if parent == nil {
			continue
		}

		parent.UsedBy = &model.UserID
		claimedAt := model.ClaimedAt
		parent.UsedAt = &claimedAt
		parent.Status = service.StatusUsed
		parent.Value = model.Amount
		out = append(out, *parent)
	}
	return out
}

func sortRedeemCodesByUsedAtDesc(codes []service.RedeemCode) {
	sort.Slice(codes, func(i, j int) bool {
		left := codes[i].UsedAt
		right := codes[j].UsedAt
		switch {
		case left == nil && right == nil:
			return codes[i].CreatedAt.After(codes[j].CreatedAt)
		case left == nil:
			return false
		case right == nil:
			return true
		case left.Equal(*right):
			return codes[i].CreatedAt.After(codes[j].CreatedAt)
		default:
			return left.After(*right)
		}
	})
}
