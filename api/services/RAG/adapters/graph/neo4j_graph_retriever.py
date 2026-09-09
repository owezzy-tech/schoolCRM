from typing import Any

from domain.entities.graph_context import GraphContext
from domain.ports.graph_retriever import IGraphRetriever


class Neo4jGraphRetriever(IGraphRetriever):
    def __init__(self, driver: Any) -> None:
        self._driver = driver

    async def retrieve(self, *, question: str, collection: str, limit: int) -> list[GraphContext]:
        query = """
        MATCH (source:AdmissionsKnowledge {collection: $collection})
            -[:HAS_CONTEXT]->(context:ContextChunk)
        WHERE toLower(context.text) CONTAINS toLower($question)
           OR toLower(source.title) CONTAINS toLower($question)
        OPTIONAL MATCH path = (source)-[:HAS_CONTEXT]->(context)
        RETURN
            coalesce(source.id, elementId(source)) AS source_id,
            labels(source)[0] AS source_type,
            source.title AS title,
            context.text AS text,
            source.collection AS collection,
            1.0 AS score,
            [node IN nodes(path) | coalesce(node.title, node.id, elementId(node))] AS path
        LIMIT $limit
        """

        result = await self._driver.execute_query(
            query,
            question=question,
            collection=collection,
            limit=limit,
        )
        return [self._to_context(record.data()) for record in result.records]

    @staticmethod
    def _to_context(data: dict[str, Any]) -> GraphContext:
        return GraphContext(
            source_id=str(data["source_id"]),
            source_type=str(data["source_type"]),
            title=str(data["title"]),
            text=str(data["text"]),
            collection=str(data["collection"]),
            score=float(data["score"]),
            path=[str(item) for item in data.get("path", [])],
        )
