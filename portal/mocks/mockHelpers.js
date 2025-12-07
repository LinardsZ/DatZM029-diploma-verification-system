/**
 * Parses request body and supports both Promise and callback patterns
 * @param {import('http').IncomingMessage} req - Request object
 * @param {((body: any) => void)} [callback] - Optional callback to handle parsed body
 * @returns {Promise<any>} Parsed JSON body or null if parsing failed
 * @example
 * ```js
 * {
 *    pattern: '/api/evidence-types',
 *    method: 'POST',
 *    handle: async (req, res) => {
 *      await handleWithBody(req, (evidenceType) => {
 *        const newItem = {
 *          ...evidenceType,
 *          evidenceTypeId: evidenceType.evidenceTypeId,
 *          name: evidenceType.name,
 *          availableTranslations: evidenceType?.availableTranslations || [],
 *          evidenceFields: evidenceType?.evidenceFields || [],
 *        };
 *
 *        evidenceTypes.push(newItem);
 *
 *        res.setHeader('Content-Type', 'application/json');
 *        res.end(JSON.stringify(newItem));
 *      });
 *    },
 *  },
 * ```
 */
export async function handleWithBody(req, callback) {
  return new Promise((resolve) => {
    let body = '';
    req.on('data', (chunk) => {
      body += chunk;
      // Check for complete JSON object
      if (body.trim().endsWith('}')) {
        try {
          const parsedBody = JSON.parse(body);
          resolve(parsedBody);
          // Optional callback support
          // eslint-disable-next-line no-unused-expressions
          callback?.call(null, parsedBody);
        } catch (error) {
          console.error('Failed to parse JSON:', error);
          resolve(null);
          // eslint-disable-next-line no-unused-expressions
          callback?.call(null, null);
        }
      }
    });
  });
}
