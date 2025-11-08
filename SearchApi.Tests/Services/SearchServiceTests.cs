using Xunit;
using Moq;
using FluentAssertions;
using SearchApi.Services;
using SearchApi.Models;
using SolrNet;
using Microsoft.Extensions.Logging;

namespace SearchApi.Tests.Services;

public class SearchServiceTests
{
    private readonly Mock<ISolrOperations<MetadataDocument>> _mockSolr;
    private readonly Mock<ILogger<SearchService>> _mockLogger;
    private readonly SearchService _service;

    public SearchServiceTests()
    {
        _mockSolr = new Mock<ISolrOperations<MetadataDocument>>();
        _mockLogger = new Mock<ILogger<SearchService>>();
        _service = new SearchService(_mockSolr.Object, _mockLogger.Object);
    }

    [Fact]
    public async Task SearchAsync_ShouldReturnResults_WhenDocumentsExist()
    {
        // Arrange
        var request = new SearchRequest
        {
            Query = "test",
            Rows = 10,
            Start = 0
        };

        var documents = new SolrQueryResults<MetadataDocument>
        {
            new MetadataDocument { Id = "1", Title = "Test Document" }
        };
        documents.NumFound = 1;

        _mockSolr.Setup(s => s.Query(It.IsAny<string>(), It.IsAny<QueryOptions>()))
            .Returns(documents);

        // Act
        var result = await _service.SearchAsync(request);

        // Assert
        result.Should().NotBeNull();
        result.TotalResults.Should().Be(1);
        result.Results.Should().HaveCount(1);
        result.Results[0].Title.Should().Be("Test Document");
    }

    [Fact]
    public async Task GetByIdAsync_ShouldReturnDocument_WhenExists()
    {
        // Arrange
        var documentId = "test-id";
        var documents = new SolrQueryResults<MetadataDocument>
        {
            new MetadataDocument { Id = documentId, Title = "Test" }
        };

        _mockSolr.Setup(s => s.Query(It.Is<string>(q => q.Contains(documentId)), null))
            .Returns(documents);

        // Act
        var result = await _service.GetByIdAsync(documentId);

        // Assert
        result.Should().NotBeNull();
        result!.Id.Should().Be(documentId);
    }

    [Fact]
    public async Task IndexDocumentAsync_ShouldReturnTrue_WhenSuccessful()
    {
        // Arrange
        var document = new MetadataDocument
        {
            Id = "test-id",
            Title = "Test Document"
        };

        _mockSolr.Setup(s => s.Add(It.IsAny<MetadataDocument>()))
            .Returns(new ResponseHeader());
        _mockSolr.Setup(s => s.Commit())
            .Returns(new ResponseHeader());

        // Act
        var result = await _service.IndexDocumentAsync(document);

        // Assert
        result.Should().BeTrue();
        _mockSolr.Verify(s => s.Add(It.IsAny<MetadataDocument>()), Times.Once);
        _mockSolr.Verify(s => s.Commit(), Times.Once);
    }

    [Fact]
    public async Task DeleteDocumentAsync_ShouldReturnTrue_WhenSuccessful()
    {
        // Arrange
        var documentId = "test-id";

        _mockSolr.Setup(s => s.Delete(It.IsAny<string>()))
            .Returns(new ResponseHeader());
        _mockSolr.Setup(s => s.Commit())
            .Returns(new ResponseHeader());

        // Act
        var result = await _service.DeleteDocumentAsync(documentId);

        // Assert
        result.Should().BeTrue();
        _mockSolr.Verify(s => s.Delete(documentId), Times.Once);
        _mockSolr.Verify(s => s.Commit(), Times.Once);
    }
}
