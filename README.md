# IdeaGen

Microservices app written on golang. Generating ideas for self-education.

## ENV FILES
- ./.env - file with main connectivity setup
    - PORT - docker port
    - IDEA_PORT - restapi port
    - DB_HOST
    - DB_PORT
    - DB_USER
    - DB_PASSWORD
    - DB_NAME
    - REDIS_HOST
    - REDIS_PORT
    - REDIS_PASSWORD

- ./app.env - file with restapi setup
    - AI_PROVIDER - provider of LLM (only "groq" supports)
    - AI_API_KEY - key of LLM provider
    - AI_MODEL - LLM model type in ai provider
    - JWT_SECRET_KEY - jwt secret
    - ACCESS_TOKEN_TTL_MINUTES
    - REFRESH_TOKEN_TTL_DAYS

- ./frontend/.env.local - file with frontend env vars
    - NEXT_PUBLIC_API_URL=/v1
