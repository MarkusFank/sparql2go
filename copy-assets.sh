# build frontend
cd sparql2go-frontend
npm run build
cd ..

# copy frontend to target folder for embedding
cp -R ./sparql2go-frontend/dist/ ./web/dist