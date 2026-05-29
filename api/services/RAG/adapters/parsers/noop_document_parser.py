from domain.ports.document_parser import IDocumentParser


class NoopDocumentParser(IDocumentParser):
    async def parse(self, *, filename: str, content_type: str, payload: bytes) -> str:
        if content_type == "text/plain":
            return payload.decode("utf-8", errors="ignore")

        return (
            f"Scaffold parser placeholder for {filename} ({content_type}). "
            "Replace with a LangChain-backed parser."
        )
