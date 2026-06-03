-- Seed default providers and models.
-- Uses ON CONFLICT DO NOTHING so it is safe to run multiple times.

INSERT INTO providers (id, name, base_url, api_key)
VALUES
    ('a0000000-0000-0000-0000-000000000001', 'OpenAI',    'https://api.openai.com',                   'sk-your-api-key'),
    ('a0000000-0000-0000-0000-000000000002', 'Anthropic', 'https://api.anthropic.com',                'sk-ant-your-api-key'),
    ('a0000000-0000-0000-0000-000000000003', 'Google',    'https://generativelanguage.googleapis.com', 'your-google-api-key')
ON CONFLICT (name) DO NOTHING;

INSERT INTO models (id, provider_id, model_name, is_enabled)
SELECT gen_random_uuid(), p.id, v.model_name, TRUE
FROM (
    VALUES
        -- OpenAI
        ('gpt-4o',      'OpenAI'),
        ('gpt-4o-mini', 'OpenAI'),
        ('gpt-4.1',     'OpenAI'),
        ('o4-mini',     'OpenAI'),
        -- Anthropic
        ('claude-sonnet-4-20250514',  'Anthropic'),
        ('claude-opus-4-20250514',    'Anthropic'),
        ('claude-haiku-4-5-20250514', 'Anthropic'),
        -- Google
        ('gemini-2.5-flash', 'Google'),
        ('gemini-2.5-pro',   'Google')
) AS v(model_name, provider_name)
JOIN providers p ON p.name = v.provider_name
ON CONFLICT (provider_id, model_name) DO NOTHING;
