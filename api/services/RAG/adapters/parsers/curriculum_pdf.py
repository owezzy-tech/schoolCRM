from io import BytesIO

from llama_index.core.node_parser import SentenceSplitter
from pypdf import PdfReader

from domain.entities.curriculum import CurriculumError, CurriculumPassage


def parse_curriculum_pdf(
    content: bytes, first_page: int, last_page: int
) -> list[CurriculumPassage]:
    """Keep PDF viewer page numbers; never infer grade boundaries or silently OCR."""
    if not content.startswith(b"%PDF-"):
        raise CurriculumError(422, "Only PDF curriculum originals are supported")
    try:
        reader = PdfReader(BytesIO(content), strict=True)
        if reader.is_encrypted or not 1 <= first_page <= last_page <= len(reader.pages):
            raise CurriculumError(422, "Invalid page range or encrypted PDF")
        splitter = SentenceSplitter(chunk_size=400, chunk_overlap=50)
        passages = []
        for page_number in range(first_page, last_page + 1):
            text = reader.pages[page_number - 1].extract_text() or ""
            if len(text.strip()) < 40:
                raise CurriculumError(422, "Selected pages require OCR and human review")
            for text_chunk in splitter.split_text(text):
                passages.append(CurriculumPassage(page_number, len(passages), text_chunk))
        return passages
    except CurriculumError:
        raise
    except Exception as exc:
        raise CurriculumError(422, "Unable to parse this PDF original") from exc
