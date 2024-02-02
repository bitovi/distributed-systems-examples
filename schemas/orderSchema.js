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
          required: ['product_code', 'quantity', 'total'],
          properties: {
            product_code: { type: 'string' },
            quantity: { type: 'number' },
            total: { type: 'number' }
          }
        }
      }
    }
  }
};

module.exports = { orderSchema };
