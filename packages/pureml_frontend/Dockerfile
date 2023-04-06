FROM node:19-alpine
WORKDIR /usr/server/app

# Install dependencies
COPY ./package.json ./
RUN npm install --legacy-peer-deps

# TODO: Build (using dev temporarily due to env variable issue)
COPY ./ .
RUN npm run build

# Set environment variables
# ENV NODE_ENV=production
ENV NEXT_PUBLIC_BACKEND_URL=http://localhost:8080/api/
ARG NEXT_PUBLIC_AIRTABLE_API_KEY
ARG GA_TRACKING_ID

# Run
# will launch the remix app when we run this Docker image.
CMD ["npm", "run" ,"dev"] 
EXPOSE 3000
