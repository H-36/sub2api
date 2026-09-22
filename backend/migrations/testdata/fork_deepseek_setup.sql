INSERT INTO accounts (id, name, platform, type, credentials, extra) VALUES
 (924001, 'fork-flash', 'openai', 'apikey', '{"api_key":"test-secret","base_url":"https://third.example/proxy?tenant=1","model_mapping":{"ignored":"old"}}', '{"provider":"DeepSeek","openai_responses_mode":"force_responses","keep":true}'),
 (924002, 'ordinary-openai', 'openai', 'apikey', '{"api_key":"ordinary"}', '{}'),
 (924003, 'native-deepseek', 'deepseek', 'apikey', '{"api_key":"native","api_protocol":"adaptive"}', '{}'),
 (924004, 'fork-v4', 'openai', 'apikey', '{"api_key":"test-v4","base_url":"https://third.example/v4/"}', '{"provider":"deepseek"}'),
 (924005, 'fork-chat', 'openai', 'apikey', '{"api_key":"chat","base_url":"https://third.example"}', '{"provider":"deepseek","openai_responses_mode":"force_chat_completions"}');
INSERT INTO groups (id, name, platform) VALUES
 (924001, 'fork-native', 'deepseek'),
 (924002, 'fork-openai-only', 'openai'),
 (924003, 'fork-mixed', 'openai'),
 (924004, 'unrelated', 'openai');
INSERT INTO account_groups (account_id, group_id) VALUES
 (924001,924001), (924001,924002), (924001,924003), (924002,924003), (924002,924004);
INSERT INTO channels (id, name, model_mapping) VALUES
 (924001, 'fork-channel', '{"openai":{"alias":"deepseek-v4-flash","kept":"original"},"deepseek":{"kept":"native"}}'),
 (924002, 'unrelated-channel', '{}');
INSERT INTO channel_groups (channel_id, group_id) VALUES (924001,924001), (924001,924003), (924002,924004);
INSERT INTO channel_model_pricing (id, channel_id, platform, models, input_price, time_pricing, reasoning_effort_multipliers) VALUES
 (924001,924001,'openai','["deepseek-v4-flash","native-priced"]',0.12,'{"timezone":"UTC","periods":[]}','{"max":2}'),
 (924002,924001,'deepseek','["NATIVE-PRICED"]',0.99,NULL,'{}'),
 (924003,924002,'openai','["unrelated"]',0.77,NULL,'{}');
INSERT INTO channel_pricing_intervals (pricing_id,min_tokens,max_tokens,input_price) VALUES (924001,100,1000,0.03);
INSERT INTO channel_account_stats_pricing_rules (id,channel_id,name) VALUES (924001,924001,'fork-costs');
INSERT INTO channel_account_stats_model_pricing (id,rule_id,platform,models,input_price,reasoning_effort_multipliers) VALUES
 (924001,924001,'openai','["deepseek-v4-flash"]',0.04,'{"high":1.5}');
INSERT INTO channel_account_stats_pricing_intervals (pricing_id,min_tokens,input_price) VALUES (924001,50,0.02);
