DO $$
BEGIN
 IF (SELECT count(*) FROM scheduler_outbox WHERE event_type='account_changed' AND account_id BETWEEN 924001 AND 924005) <> 3
    OR (SELECT count(*) FROM scheduler_outbox WHERE event_type='group_changed' AND group_id BETWEEN 924001 AND 924004) <> 3
    THEN RAISE EXCEPTION 'scheduler invalidation missing or replayed'; END IF;
 IF NOT EXISTS (SELECT 1 FROM accounts WHERE id=924001 AND platform='deepseek'
     AND credentials->>'api_key'='test-secret' AND credentials->>'base_url'='https://third.example/proxy/v1?tenant=1'
     AND credentials->>'api_protocol'='responses'
     AND credentials->'model_mapping'='{"deepseek-v4-flash":"deepseek-v4-flash"}'::jsonb
     AND extra='{"keep":true}'::jsonb) THEN RAISE EXCEPTION 'legacy account migration failed'; END IF;
 IF NOT EXISTS (SELECT 1 FROM accounts WHERE id=924004 AND credentials->>'base_url'='https://third.example/v4')
    OR NOT EXISTS (SELECT 1 FROM accounts WHERE id=924005 AND credentials->>'base_url'='https://third.example' AND credentials->>'api_protocol'='chat_completions')
    THEN RAISE EXCEPTION 'protocol endpoint changed'; END IF;
 IF NOT EXISTS (SELECT 1 FROM accounts WHERE id=924002 AND platform='openai' AND credentials='{"api_key":"ordinary"}'::jsonb)
    OR NOT EXISTS (SELECT 1 FROM accounts WHERE id=924003 AND credentials='{"api_key":"native","api_protocol":"adaptive"}'::jsonb)
    THEN RAISE EXCEPTION 'unrelated/native account changed'; END IF;
 IF (SELECT count(*) FROM account_groups WHERE account_id=924001) <> 3
    OR (SELECT platform FROM groups WHERE id=924001) <> 'deepseek'
    OR (SELECT platform FROM groups WHERE id=924002) <> 'deepseek'
    OR (SELECT platform FROM groups WHERE id=924003) <> 'composite'
    OR (SELECT platform FROM groups WHERE id=924004) <> 'openai'
    THEN RAISE EXCEPTION 'group or membership migration failed'; END IF;
 IF (SELECT count(*) FROM composite_model_routes WHERE group_id=924003 AND public_model='deepseek-v4-flash' AND target_platform='deepseek') <> 1
    THEN RAISE EXCEPTION 'mixed pool routing failed'; END IF;
 IF (SELECT model_mapping->'deepseek' FROM channels WHERE id=924001) <> '{"alias":"deepseek-v4-flash","kept":"native"}'::jsonb
    OR (SELECT model_mapping->'openai'->>'kept' FROM channels WHERE id=924001) <> 'original'
    THEN RAISE EXCEPTION 'channel mappings changed'; END IF;
 IF (SELECT count(*) FROM channel_model_pricing WHERE channel_id=924001) <> 3
    OR (SELECT count(*) FROM channel_model_pricing WHERE channel_id=924002) <> 1
    THEN RAISE EXCEPTION 'duplicate/missing/unrelated pricing'; END IF;
 IF NOT EXISTS (SELECT 1 FROM channel_model_pricing p JOIN channel_pricing_intervals i ON i.pricing_id=p.id
     WHERE p.channel_id=924001 AND p.platform='deepseek' AND p.models='["deepseek-v4-flash"]'::jsonb
       AND p.input_price=0.12 AND p.reasoning_effort_multipliers='{"max":2}'::jsonb
       AND p.time_pricing='{"timezone":"UTC","periods":[]}'::jsonb AND i.input_price=0.03 AND i.min_tokens=100 AND i.max_tokens=1000)
    THEN RAISE EXCEPTION 'pricing fields/tiers lost'; END IF;
 IF NOT EXISTS (SELECT 1 FROM channel_account_stats_model_pricing p JOIN channel_account_stats_pricing_intervals i ON i.pricing_id=p.id
     WHERE p.rule_id=924001 AND p.platform='deepseek' AND p.input_price=0.04
       AND p.reasoning_effort_multipliers='{"high":1.5}'::jsonb AND i.min_tokens=50 AND i.input_price=0.02)
    OR (SELECT count(*) FROM channel_account_stats_model_pricing WHERE rule_id=924001) <> 2
    THEN RAISE EXCEPTION 'account cost pricing lost/duplicated'; END IF;
END $$;
