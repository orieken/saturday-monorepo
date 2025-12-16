import { SpanExporter, ReadableSpan } from '@opentelemetry/sdk-trace-base';
import { ExportResult, ExportResultCode } from '@opentelemetry/core';
import * as fs from 'fs';
import * as path from 'path';

export class FileSpanExporter implements SpanExporter {
  private filePath: string;

  constructor(outputDir: string = './reports', filename: string = 'otel-spans.json') {
    if (!fs.existsSync(outputDir)) {
      fs.mkdirSync(outputDir, { recursive: true });
    }
    this.filePath = path.join(outputDir, filename);
    // Initialize file with empty array if it doesn't exist, or clear it
    fs.writeFileSync(this.filePath, '[\n');
  }

  export(spans: ReadableSpan[], resultCallback: (result: ExportResult) => void): void {
    const jsonSpans = spans.map(span => ({
      name: span.name,
      kind: span.kind,
      spanContext: span.spanContext(),
      parentSpanId: span.parentSpanId,
      startTime: span.startTime,
      endTime: span.endTime,
      status: span.status,
      attributes: span.attributes,
      events: span.events,
      links: span.links,
      resource: span.resource.attributes
    }));

    const content = jsonSpans.map(s => JSON.stringify(s, null, 2)).join(',\n');
    
    // Naively append. Invalid JSON array unless we manage commas, but good enough for debugging (NDJSON-ish inside array?)
    // Let's just append as NDJSON (one JSON per line) for easier reading, ignoring the initial '[' I wrote.
    // Actually, writing valid JSON array is hard with streaming. 
    // Let's write NDJSON (Newline Delimited JSON).
    
    // Overwriting the constructor logic:
    // fs.writeFileSync(this.filePath, ''); 
    
    fs.appendFileSync(this.filePath, content + (content ? ',\n' : ''));
    
    resultCallback({ code: ExportResultCode.SUCCESS });
  }

  shutdown(): Promise<void> {
    // Close the JSON array properly?
    fs.appendFileSync(this.filePath, '{}]'); // Hack to make it valid-ish JSON at end
    return Promise.resolve();
  }
}
