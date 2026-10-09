// Local OCR for the supplied scanned table. Rebuilding JSON graphs needs only Python.
import Foundation
import PDFKit
import Vision
import AppKit
import CryptoKit

guard CommandLine.arguments.count == 3 else {
    fputs("Usage: swift scripts/ocr_rag_cutoffs.swift INPUT.pdf OUTPUT.json\n", stderr)
    exit(1)
}
let input = URL(fileURLWithPath: CommandLine.arguments[1])
guard let document = PDFDocument(url: input), document.pageCount == 4 else {
    fputs("Expected the four-page Maseno cutoff document\n", stderr)
    exit(1)
}
var pages: [[String: Any]] = []
for index in 0..<document.pageCount {
    let page = document.page(at: index)!
    let image = page.thumbnail(of: NSSize(width: 2400, height: 3400), for: .mediaBox)
    let cgImage = image.cgImage(forProposedRect: nil, context: nil, hints: nil)!
    let request = VNRecognizeTextRequest()
    request.recognitionLevel = .accurate
    request.recognitionLanguages = ["en-US"]
    request.usesLanguageCorrection = false
    try VNImageRequestHandler(cgImage: cgImage).perform([request])
    let lines: [[String: Any]] = (request.results ?? []).compactMap { observation in
        guard let candidate = observation.topCandidates(1).first else { return nil }
        return ["text": candidate.string, "confidence": candidate.confidence,
                "x": observation.boundingBox.minX, "y": observation.boundingBox.minY,
                "width": observation.boundingBox.width, "height": observation.boundingBox.height]
    }
    pages.append(["page": index + 1, "lines": lines])
}
let checksum = SHA256.hash(data: try Data(contentsOf: input)).map { String(format: "%02x", $0) }.joined()
let output: [String: Any] = ["source_sha256": checksum,
    "engine": "Apple Vision accurate en-US; language correction disabled", "pages": pages]
let data = try JSONSerialization.data(withJSONObject: output, options: [.prettyPrinted, .sortedKeys])
try data.write(to: URL(fileURLWithPath: CommandLine.arguments[2]), options: .atomic)
