const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api/v1'

async function request(path, options = {}, token) {
  const headers = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  }
  if (token) headers.Authorization = `Bearer ${token}`

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  })

  const body = await response.json().catch(() => ({}))
  if (!response.ok || body.success === false) {
    const message = body?.error?.message || 'Terjadi kesalahan'
    const details = body?.error?.details || null
    throw { message, details }
  }

  return body
}

export async function login(username, password) {
  return request('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })
}

export async function fetchProducts(params, token) {
  const query = new URLSearchParams(params)
  return request(`/products?${query.toString()}`, { method: 'GET' }, token)
}

export async function createProduct(payload, token) {
  return request('/products', {
    method: 'POST',
    body: JSON.stringify(payload),
  }, token)
}

export async function updateProduct(id, payload, token) {
  return request(`/products/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  }, token)
}

export async function deleteProduct(id, token) {
  return request(`/products/${id}`, {
    method: 'DELETE',
  }, token)
}
