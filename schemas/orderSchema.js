const orderSchema = {
  body: {
    type: 'object',
    required: ['customer', 'total', 'products'],
    properties: {
      total: { type: 'number' },
      customer: { 
        type: 'object',
        required: ['firstName', 'lastName', 'city'],
        properties: {
          firstName: { type: 'string' },
          lastName: { type: 'string' },
          city: { type: 'string' }
        }
      },
      products: { 
        type: 'array',
        minItems: 1,
        items: {
          type: 'object',
          required: ['productCode', 'quantity', 'total'],
          properties: {
            productCode: { type: 'string' },
            quantity: { type: 'number' },
            total: { type: 'number' }
          }
        }
      }
    }
  }
};

module.exports = { orderSchema };
