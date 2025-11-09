using Serilog;
using SearchApi.Services;
using SearchApi.Models;
using SolrNet;

var builder = WebApplication.CreateBuilder(args);

// Configure Serilog
Log.Logger = new LoggerConfiguration()
    .ReadFrom.Configuration(builder.Configuration)
    .Enrich.FromLogContext()
    .WriteTo.Console(formatProvider: System.Globalization.CultureInfo.InvariantCulture)
    .CreateLogger();

builder.Host.UseSerilog();

// Add services to the container
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1", new()
    {
        Title = "Search API",
        Version = "v1",
        Description = "Metadata search API using Solr"
    });
});

// Configure Solr
var solrUrl = builder.Configuration["Solr:Url"] ?? "http://solr:8983/solr/metadata";
builder.Services.AddSolrNet<MetadataDocument>(solrUrl);
builder.Services.AddScoped<ISearchService, SearchService>();

// Add health checks
builder.Services.AddHealthChecks();
// Note: Solr health check has API compatibility issues, skipping for now
// TODO: Fix Solr health check configuration once API is clarified

var app = builder.Build();

// Configure the HTTP request pipeline
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();
app.UseAuthorization();
app.MapControllers();
app.MapHealthChecks("/health");

Log.Information("Starting Search API");
app.Run();

public partial class Program { } // For testing
