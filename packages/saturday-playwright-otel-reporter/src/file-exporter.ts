import { SpanExporter, ReadableSpan } from '@opentelemetry/sdk-trace-base';
import { ExportResult, ExportResultCode } from '@opentelemetry/core';
import * as fs from 'fs';
import * as path from 'path';

export class FileSpanExporter implements SpanExporter {
  private filePath: string;

  constructor(outputDir: string = './reports', filename: string = 'otel-playwright-spans.json') {
    if (!fs.existsSync(outputDir)) {
      fs.mkdirSync(outputDir, { recursive: true });
    }
    this.filePath = path.join(outputDir, filename);
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
    fs.appendFileSync(this.filePath, content + (content ? ',\n' : ''));
    
    resultCallback({ code: ExportResultCode.SUCCESS });
  }

  shutdown(): Promise<void> {
    fs.appendFileSync(this.filePath, '{}]');
    return Promise.resolve();
  }
}
