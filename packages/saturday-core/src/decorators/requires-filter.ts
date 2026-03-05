export class FilterError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'FilterError';
  }
}

export function RequiresFilter(conditionMethod: string) {
  return function (target: any, propertyKey: string) {
    const backingField = Symbol(`__${propertyKey}`);

    Object.defineProperty(target, propertyKey, {
      configurable: true,
      enumerable: true,
      set(value: any) {
        this[backingField] = createFilteredProxy(value, this, conditionMethod, propertyKey);
      },
      get() {
        return this[backingField];
      },
    });
  };
}

function createFilteredProxy(
  element: any,
  context: any,
  conditionMethod: string,
  elementName: string
): any {
  if (!element) return element;

  return new Proxy(element, {
    get(target, prop, receiver) {
      const originalValue = Reflect.get(target, prop, receiver);

      if (typeof originalValue === 'function') {
        return async function (this: any, ...args: any[]) {
          if (typeof context[conditionMethod] !== 'function') {
            throw new Error(
              `Filter condition method '${conditionMethod}' not found on ${context.constructor.name}`
            );
          }

          const isAllowed = await context[conditionMethod]();
          
          if (!isAllowed) {
            throw new FilterError(
              `Filter check failed: '${conditionMethod}' returned false for element '${elementName}'`
            );
          }

          return originalValue.apply(this, args);
        };
      }

      return originalValue;
    },
  });
}
