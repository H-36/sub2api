# Retiring the fork DeepSeek adapter

The runtime now uses upstream DeepSeek accounts, group routing, model mappings,
pricing and UI. There is no fork provider flag, OpenAI-to-DeepSeek group alias,
hard-coded Flash restriction or separate DeepSeek icon.

Migration `240_migrate_fork_deepseek_accounts.sql` runs once during the normal
database migration phase on deployment. It does not contact provider APIs.

- Only `platform=openai`, `type=apikey`, `extra.provider=deepseek` accounts are
  converted. Existing upstream DeepSeek and ordinary OpenAI accounts are untouched.
- Account IDs, credentials, proxies, account/group memberships, API key group IDs,
  subscription bindings and usage history are retained.
- The old effective Flash-only model selection becomes an editable upstream
  `model_mapping` entry. Other models can subsequently be enabled in the normal
  account editor. Previously ignored old model mappings are replaced deliberately.
- The third-party endpoint and its query string are retained. Responses URLs get
  an explicit `/v1` where the old OpenAI transport implicitly added it. The old
  Responses/Chat Completions protocol is pinned; it is not changed to adaptive
  routing or redirected to an official provider endpoint.
- OpenAI/DeepSeek groups containing only converted/native DeepSeek accounts become
  native DeepSeek groups. Mixed pools become upstream composite groups and get an
  explicit Flash-to-DeepSeek route. They use upstream model ownership/detection for
  other models; check custom aliases and add explicit routes in the group editor
  if those aliases have no discoverable account mapping. Existing composite group
  routes are left authoritative and should be reviewed if they explicitly direct
  DeepSeek traffic to OpenAI accounts.
- Shared channel OpenAI pricing is retained and copied to DeepSeek, including tiers,
  time pricing and account-cost pricing. Existing native DeepSeek entries/mappings
  take precedence. No historical billing records are rewritten.
- Scheduler invalidation events refresh the shared snapshots. Replaying the
  migration does nothing once the legacy provider markers have been removed.

Back up the database before deployment. Reverting code alone does not reverse a
data migration; use the pre-deployment database backup for a complete rollback.

The fixture in `backend/migrations/testdata/fork_deepseek_*` tests endpoint and
credential preservation, pure/mixed pools, native configuration precedence,
pricing tiers and replay. It is also exercised by `TestForkDeepSeekMigration` in
the repository integration suite.
