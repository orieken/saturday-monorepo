const mockTest: any = {
  extend: jest.fn((fixtures) => {
    mockTest.fixtures = fixtures;
    return mockTest;
  }),
  expect: jest.fn(),
};

export const test = mockTest;

export const request = {
  newContext: jest.fn(),
};

export const APIRequestContext = jest.fn();

