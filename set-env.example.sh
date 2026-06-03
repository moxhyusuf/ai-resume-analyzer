export PORT=8080
export DB_USER=resume_user
export DB_PASSWORD=<PASSWORD>
export DB_NAME=ai_resume_analyzer
export DATABASE_URL="host=postgres user=resume_user password=<PASSWORD> dbname=ai_resume_analyzer port=5432 sslmode=require"
export JWT_SECRET=$(openssl rand -base64 32)
export GROQ_API_KEY=<GROQ_API_KEY>