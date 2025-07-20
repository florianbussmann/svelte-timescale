FROM node:22-alpine AS base
WORKDIR /app
RUN corepack enable

FROM base AS prod
COPY pnpm-lock.yaml /app
RUN pnpm fetch --prod
COPY . /app
RUN pnpm run build

FROM base
COPY --from=prod /app/node_modules /app/node_modules
COPY --from=prod /app/build /app/build
EXPOSE 3000
CMD [ "node", "build" ]