import { describe, it, expect, vi, beforeEach, beforeAll } from 'vitest';

const { extendMock } = vi.hoisted(() => {
  return { extendMock: vi.fn() };
});

vi.mock('@playwright/test', () => {
  return {
    test: {
      extend: extendMock,
    },
    expect: vi.fn(),
  };
});

import { ConsoleLogger, TabManager, SiteManager } from '@orieken/saturday-core';

vi.mock('@orieken/saturday-core', () => {
  const ConsoleLoggerMock = vi.fn().mockImplementation(function() {
    return { getFormattedLogs: vi.fn().mockReturnValue('mocked logs') };
  });
  const TabManagerMock = vi.fn().mockImplementation(function() {
    return {
      getActiveName: vi.fn().mockReturnValue('mockTab'),
      count: vi.fn().mockReturnValue(1),
      closeAll: vi.fn().mockResolvedValue(undefined),
    };
  });
  const SiteManagerMock = vi.fn().mockImplementation(function() {
    return {
      getActiveName: vi.fn().mockReturnValue('mockSite'),
      listSites: vi.fn().mockReturnValue(['mockSite']),
      clear: vi.fn(),
    };
  });

  return {
    ConsoleLogger: ConsoleLoggerMock,
    TabManager: TabManagerMock,
    SiteManager: SiteManagerMock,
  };
});

// Import after mocks are set up
import '../src/index';

describe('saturday-playwright fixtures', () => {
  let fixtures: any;

  beforeAll(() => {
    // The first call to base.extend will have our fixtures object.
    fixtures = extendMock.mock.calls[0][0];
  });

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('consoleLogger fixture', () => {
    it('sets up ConsoleLogger and calls use', async () => {
      const useMock = vi.fn();
      const mockPage = {};
      const mockTestInfo = { status: 'passed', expectedStatus: 'passed' };

      await fixtures.consoleLogger({ page: mockPage }, useMock, mockTestInfo);

      expect(ConsoleLogger).toHaveBeenCalledWith(mockPage);
      expect(useMock).toHaveBeenCalled();
    });

    it('attaches logs on failure and logs exist', async () => {
      const useMock = vi.fn();
      const mockPage = {};
      const attachMock = vi.fn();
      const mockTestInfo = { status: 'failed', expectedStatus: 'passed', attach: attachMock };

      await fixtures.consoleLogger({ page: mockPage }, useMock, mockTestInfo);

      expect(attachMock).toHaveBeenCalledWith('Console Logs', {
        contentType: 'text/plain',
        body: expect.any(Buffer),
      });
    });

    it('does not attach logs if none exist', async () => {
      // Override mock just for this test
      vi.mocked(ConsoleLogger).mockImplementationOnce(function() {
        return { getFormattedLogs: vi.fn().mockReturnValue('') };
      } as any);

      const useMock = vi.fn();
      const mockPage = {};
      const attachMock = vi.fn();
      const mockTestInfo = { status: 'failed', expectedStatus: 'passed', attach: attachMock };

      await fixtures.consoleLogger({ page: mockPage }, useMock, mockTestInfo);

      expect(attachMock).not.toHaveBeenCalled();
    });
  });

  describe('tabManager fixture', () => {
    it('sets up TabManager and calls use and cleanup', async () => {
      const useMock = vi.fn();
      const mockContext = {};
      const mockPage = {};
      const mockTestInfo = { status: 'passed', expectedStatus: 'passed' };

      await fixtures.tabManager({ context: mockContext, page: mockPage }, useMock, mockTestInfo);

      expect(TabManager).toHaveBeenCalledWith(mockContext, mockPage);
      expect(useMock).toHaveBeenCalled();
      
      // Get the mocked instance
      const mockedInstance = vi.mocked(TabManager).mock.results[0].value;
      expect(mockedInstance.closeAll).toHaveBeenCalledWith(true);
    });

    it('attaches tab state on failure', async () => {
      const useMock = vi.fn();
      const mockContext = {};
      const mockPage = {};
      const attachMock = vi.fn();
      const mockTestInfo = { status: 'failed', expectedStatus: 'passed', attach: attachMock };

      await fixtures.tabManager({ context: mockContext, page: mockPage }, useMock, mockTestInfo);

      expect(attachMock).toHaveBeenCalledWith('Active Tab State', {
        contentType: 'application/json',
        body: expect.any(Buffer),
      });
    });
  });

  describe('siteManager fixture', () => {
    it('sets up SiteManager and calls use and cleanup', async () => {
      const useMock = vi.fn();
      const mockPage = {};
      const mockTestInfo = { status: 'passed', expectedStatus: 'passed' };

      await fixtures.siteManager({ page: mockPage }, useMock, mockTestInfo);

      expect(SiteManager).toHaveBeenCalledWith(mockPage);
      expect(useMock).toHaveBeenCalled();
      
      const mockedInstance = vi.mocked(SiteManager).mock.results[0].value;
      expect(mockedInstance.clear).toHaveBeenCalled();
    });

    it('attaches site state on failure', async () => {
      const useMock = vi.fn();
      const mockPage = {};
      const attachMock = vi.fn();
      const mockTestInfo = { status: 'failed', expectedStatus: 'passed', attach: attachMock };

      await fixtures.siteManager({ page: mockPage }, useMock, mockTestInfo);

      expect(attachMock).toHaveBeenCalledWith('Site Manager State', {
        contentType: 'application/json',
        body: expect.any(Buffer),
      });
    });
  });
});
