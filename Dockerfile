# Build stage
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

# Copy solution and project files
COPY SearchApi.sln .
COPY SearchApi/SearchApi.csproj SearchApi/
COPY SearchApi.Tests/SearchApi.Tests.csproj SearchApi.Tests/
COPY SearchApi.IntegrationTests/SearchApi.IntegrationTests.csproj SearchApi.IntegrationTests/

# Restore dependencies
RUN dotnet restore SearchApi.sln

# Copy source code
COPY SearchApi/ SearchApi/
COPY SearchApi.Tests/ SearchApi.Tests/
COPY SearchApi.IntegrationTests/ SearchApi.IntegrationTests/

# Build the application
WORKDIR /src/SearchApi
RUN dotnet build -c Release -o /app/build --no-restore

# Run tests
WORKDIR /src
RUN dotnet test SearchApi.Tests/SearchApi.Tests.csproj -c Release --no-build --verbosity normal

# Publish stage
FROM build AS publish
WORKDIR /src/SearchApi
RUN dotnet publish -c Release -o /app/publish --no-restore /p:UseAppHost=false

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime

# Create non-root user for security
RUN groupadd -r searchapi && useradd -r -g searchapi searchapi

WORKDIR /app

# Copy published files
COPY --from=publish /app/publish .

# Set ownership
RUN chown -R searchapi:searchapi /app

# Switch to non-root user
USER searchapi

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Set environment variables
ENV ASPNETCORE_URLS=http://+:8080
ENV DOTNET_RUNNING_IN_CONTAINER=true
ENV DOTNET_SYSTEM_GLOBALIZATION_INVARIANT=false

ENTRYPOINT ["dotnet", "SearchApi.dll"]
