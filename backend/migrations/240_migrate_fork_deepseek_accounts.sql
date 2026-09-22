-- Retire the fork's openai + extra.provider=deepseek representation.
-- This is a data conversion only: runtime routing uses upstream platforms.
-- IDs, API keys, proxies, group memberships and historical usage stay intact.
SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '10min';

CREATE TEMP TABLE fork_deepseek_accounts_240 ON COMMIT DROP AS
SELECT id FROM accounts
WHERE platform = 'openai' AND type = 'apikey' AND deleted_at IS NULL
  AND lower(btrim(extra->>'provider')) = 'deepseek';

CREATE TEMP TABLE fork_deepseek_groups_240 ON COMMIT DROP AS
SELECT DISTINCT g.id, g.platform AS original_platform
FROM groups g JOIN account_groups ag ON ag.group_id = g.id
JOIN fork_deepseek_accounts_240 old ON old.id = ag.account_id;

CREATE TEMP TABLE fork_deepseek_channels_240 ON COMMIT DROP AS
SELECT DISTINCT cg.channel_id FROM channel_groups cg
JOIN fork_deepseek_groups_240 g ON g.id = cg.group_id;

-- Keep the old third-party Responses URL: OpenAI adds /v1 by default,
-- whereas native DeepSeek adds /responses directly. Preserve query strings.
WITH old_urls AS (
    SELECT a.id,
           regexp_replace(COALESCE(NULLIF(btrim(a.credentials->>'base_url'), ''),
                                   'https://api.openai.com'), '#.*$', '') AS base
    FROM accounts a JOIN fork_deepseek_accounts_240 old ON old.id = a.id
), urls AS (
    SELECT id, rtrim(split_part(base, '?', 1), '/') AS path,
           COALESCE(substring(base FROM '\?.*$'), '') AS query FROM old_urls
)
UPDATE accounts a
SET platform = 'deepseek',
    credentials = COALESCE(a.credentials, '{}'::jsonb) || jsonb_build_object(
        'account_mode', 'payg',
        'api_protocol', CASE WHEN a.extra->>'openai_responses_mode' = 'force_chat_completions'
                             THEN 'chat_completions' ELSE 'responses' END,
        'base_url', CASE
            WHEN a.extra->>'openai_responses_mode' = 'force_chat_completions'
                 OR urls.path ~* '/(responses|v[0-9]+(\.[0-9]+|alpha[^/]*|beta[^/]*|preview[^/]*)?)$'
            THEN urls.path || urls.query
            ELSE urls.path || '/v1' || urls.query END,
        -- Replace the old hard-coded Flash restriction with editable upstream config.
        'model_mapping', jsonb_build_object('deepseek-v4-flash', 'deepseek-v4-flash')),
    extra = COALESCE(a.extra, '{}'::jsonb) - 'provider' - 'openai_responses_mode'
            - 'openai_passthrough' - 'openai_oauth_passthrough',
    updated_at = NOW()
FROM urls WHERE a.id = urls.id;

-- Preserve existing API key/group IDs. Mixed pools become upstream composite
-- groups so ordinary OpenAI accounts remain reachable alongside DeepSeek.
UPDATE groups g
SET platform = CASE WHEN EXISTS (
        SELECT 1 FROM account_groups ag JOIN accounts a ON a.id = ag.account_id
        WHERE ag.group_id = g.id AND a.deleted_at IS NULL AND a.platform <> 'deepseek'
    ) THEN 'composite' ELSE 'deepseek' END,
    updated_at = NOW()
FROM fork_deepseek_groups_240 old
WHERE g.id = old.id AND old.original_platform IN ('openai', 'deepseek');

-- Existing explicit routes are authoritative. For newly converted mixed pools,
-- disambiguate Flash (ordinary OpenAI accounts may expose it as well).
INSERT INTO composite_model_routes (group_id, public_model, target_platform, notes)
SELECT g.id, 'deepseek-v4-flash', 'deepseek', 'Migrated fork DeepSeek account pool'
FROM groups g JOIN fork_deepseek_groups_240 old ON old.id = g.id
WHERE g.platform = 'composite' AND old.original_platform <> 'composite'
ON CONFLICT (group_id, endpoint, match_type, public_model)
WHERE deleted_at IS NULL DO NOTHING;

-- DeepSeek groups previously read OpenAI channel mappings. Copy them without
-- overwriting a native DeepSeek mapping or changing the shared OpenAI mapping.
UPDATE channels c
SET model_mapping = jsonb_set(COALESCE(c.model_mapping, '{}'::jsonb), '{deepseek}',
        c.model_mapping->'openai' || CASE WHEN jsonb_typeof(c.model_mapping->'deepseek') = 'object'
            THEN c.model_mapping->'deepseek' ELSE '{}'::jsonb END),
    updated_at = NOW()
WHERE c.id IN (SELECT channel_id FROM fork_deepseek_channels_240)
  AND jsonb_typeof(c.model_mapping->'openai') = 'object';

-- Copy pricing and all its tiers, including time/effort/image-token pricing.
-- Populate records by column name so no optional price field is lost. Existing
-- native rows win; retain OpenAI rows for other groups sharing the channel.
DO $$
DECLARE
    pricing_table text;
    interval_table text;
    owner_column text;
    eligibility text;
    source_row record;
    new_id bigint;
    missing_models jsonb;
BEGIN
    FOR pricing_table, interval_table, owner_column, eligibility IN
        SELECT * FROM (VALUES
            ('channel_model_pricing', 'channel_pricing_intervals', 'channel_id',
             'p.channel_id IN (SELECT channel_id FROM fork_deepseek_channels_240)'),
            ('channel_account_stats_model_pricing', 'channel_account_stats_pricing_intervals', 'rule_id',
             'p.rule_id IN (SELECT id FROM channel_account_stats_pricing_rules WHERE channel_id IN (SELECT channel_id FROM fork_deepseek_channels_240))')
        ) AS tables(pricing, intervals, owner, eligible)
    LOOP
        FOR source_row IN EXECUTE format(
            'SELECT p.id, p.%I AS owner_id, to_jsonb(p) AS data FROM %I p WHERE p.platform = ''openai'' AND %s ORDER BY p.id',
            owner_column, pricing_table, eligibility)
        LOOP
            EXECUTE format(
                'SELECT jsonb_agg(m.model ORDER BY m.ordinal) FROM jsonb_array_elements($1->''models'') WITH ORDINALITY m(model, ordinal)
                 WHERE NOT EXISTS (SELECT 1 FROM %I native, jsonb_array_elements_text(native.models) existing(model)
                     WHERE native.%I = $2 AND native.platform = ''deepseek''
                       AND lower(btrim(existing.model)) = lower(btrim(m.model #>> ''{}'')))',
                pricing_table, owner_column)
                INTO missing_models USING source_row.data, source_row.owner_id;
            IF missing_models IS NULL THEN CONTINUE; END IF;
            new_id := nextval(pg_get_serial_sequence(pricing_table, 'id'));
            EXECUTE format('INSERT INTO %I SELECT (jsonb_populate_record(NULL::%I, $1)).*', pricing_table, pricing_table)
                USING source_row.data || jsonb_build_object('id', new_id, 'platform', 'deepseek', 'models', missing_models);
            EXECUTE format(
                'INSERT INTO %I SELECT (jsonb_populate_record(NULL::%I, to_jsonb(tier) || jsonb_build_object(''id'', nextval(pg_get_serial_sequence(%L, ''id'')), ''pricing_id'', $1))).*
                 FROM %I tier WHERE tier.pricing_id = $2',
                interval_table, interval_table, interval_table, interval_table)
                USING new_id, source_row.id;
        END LOOP;
    END LOOP;
END $$;

-- Refresh shared scheduler snapshots after the account/group platform changes.
INSERT INTO scheduler_outbox (event_type, account_id)
SELECT 'account_changed', id FROM fork_deepseek_accounts_240;
INSERT INTO scheduler_outbox (event_type, group_id)
SELECT 'group_changed', id FROM fork_deepseek_groups_240;

DROP TABLE fork_deepseek_channels_240;
DROP TABLE fork_deepseek_groups_240;
DROP TABLE fork_deepseek_accounts_240;
