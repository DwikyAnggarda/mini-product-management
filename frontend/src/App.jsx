import { useEffect, useMemo, useState } from 'react'
import { createProduct, deleteProduct, fetchProducts, login, updateProduct } from './api'
import './styles.css'

const emptyForm = {
  sku: '',
  name: '',
  description: '',
  price: 0,
  status: 'active',
}

function App() {
  const [token, setToken] = useState(() => localStorage.getItem('pm_token') || '')
  const [authForm, setAuthForm] = useState({ username: 'admin', password: 'password' })
  const [items, setItems] = useState([])
  const [meta, setMeta] = useState({ page: 1, limit: 10, total: 0, total_pages: 1 })
  const [filters, setFilters] = useState({ q: '', status: '', page: 1, limit: 10 })
  const [form, setForm] = useState(emptyForm)
  const [editingId, setEditingId] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [validation, setValidation] = useState({})

  const isAuthed = useMemo(() => Boolean(token), [token])

  async function loadProducts(nextFilters = filters) {
    if (!token) return
    setLoading(true)
    setError('')
    try {
      const result = await fetchProducts(nextFilters, token)
      setItems(result.data || [])
      setMeta(result.meta || { page: 1, limit: 10, total: 0, total_pages: 1 })
    } catch (err) {
      setError(err.message || 'Gagal load data produk')
      if (err.message?.toLowerCase().includes('token')) {
        handleLogout()
      }
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadProducts()
  }, [token])

  async function handleLogin(event) {
    event.preventDefault()
    setError('')
    try {
      const result = await login(authForm.username, authForm.password)
      const nextToken = result.data.token
      localStorage.setItem('pm_token', nextToken)
      setToken(nextToken)
    } catch (err) {
      setError(err.message || 'Login gagal')
    }
  }

  function handleLogout() {
    localStorage.removeItem('pm_token')
    setToken('')
    setItems([])
  }

  async function handleSearch(event) {
    event.preventDefault()
    const next = { ...filters, page: 1 }
    setFilters(next)
    await loadProducts(next)
  }

  async function handleSubmit(event) {
    event.preventDefault()
    setError('')
    setValidation({})

    const payload = {
      ...form,
      price: Number(form.price),
    }

    try {
      if (editingId) {
        await updateProduct(editingId, payload, token)
      } else {
        await createProduct(payload, token)
      }
      setForm(emptyForm)
      setEditingId(null)
      await loadProducts(filters)
    } catch (err) {
      setError(err.message || 'Gagal simpan produk')
      if (err.details) setValidation(err.details)
    }
  }

  async function handleDelete(id) {
    const ok = window.confirm('Yakin ingin menghapus produk ini?')
    if (!ok) return

    try {
      await deleteProduct(id, token)
      await loadProducts(filters)
    } catch (err) {
      setError(err.message || 'Gagal menghapus produk')
    }
  }

  function handleEdit(item) {
    setEditingId(item.id)
    setForm({
      sku: item.sku,
      name: item.name,
      description: item.description,
      price: item.price,
      status: item.status,
    })
    setValidation({})
  }

  async function goToPage(page) {
    const next = { ...filters, page }
    setFilters(next)
    await loadProducts(next)
  }

  if (!isAuthed) {
    return (
      <main className="container auth-box">
        <h1>Product Management</h1>
        <p>Silakan login untuk mengakses dashboard.</p>
        {error && <div className="error">{error}</div>}
        <form onSubmit={handleLogin} className="card form-grid">
          <label>
            Username
            <input
              value={authForm.username}
              onChange={(e) => setAuthForm((prev) => ({ ...prev, username: e.target.value }))}
            />
          </label>
          <label>
            Password
            <input
              type="password"
              value={authForm.password}
              onChange={(e) => setAuthForm((prev) => ({ ...prev, password: e.target.value }))}
            />
          </label>
          <button type="submit">Login</button>
        </form>
      </main>
    )
  }

  return (
    <main className="container">
      <header className="topbar">
        <h1>Product Management</h1>
        <button onClick={handleLogout} className="secondary">Logout</button>
      </header>

      {error && <div className="error">{error}</div>}

      <section className="card">
        <h2>{editingId ? 'Edit Produk' : 'Tambah Produk'}</h2>
        <form onSubmit={handleSubmit} className="form-grid cols-2">
          <label>
            SKU
            <input value={form.sku} onChange={(e) => setForm((p) => ({ ...p, sku: e.target.value }))} />
            {validation.sku && <small className="field-error">{validation.sku}</small>}
          </label>
          <label>
            Nama
            <input value={form.name} onChange={(e) => setForm((p) => ({ ...p, name: e.target.value }))} />
            {validation.name && <small className="field-error">{validation.name}</small>}
          </label>
          <label>
            Deskripsi
            <textarea value={form.description} onChange={(e) => setForm((p) => ({ ...p, description: e.target.value }))} />
          </label>
          <label>
            Harga
            <input
              type="number"
              min="0"
              step="0.01"
              value={form.price}
              onChange={(e) => setForm((p) => ({ ...p, price: e.target.value }))}
            />
            {validation.price && <small className="field-error">{validation.price}</small>}
          </label>
          <label>
            Status
            <select value={form.status} onChange={(e) => setForm((p) => ({ ...p, status: e.target.value }))}>
              <option value="active">active</option>
              <option value="inactive">inactive</option>
            </select>
            {validation.status && <small className="field-error">{validation.status}</small>}
          </label>
          <div className="actions">
            <button type="submit">{editingId ? 'Update' : 'Tambah'}</button>
            {editingId && (
              <button type="button" className="secondary" onClick={() => {
                setEditingId(null)
                setForm(emptyForm)
              }}>
                Batal
              </button>
            )}
          </div>
        </form>
      </section>

      <section className="card">
        <h2>Daftar Produk</h2>
        <form className="toolbar" onSubmit={handleSearch}>
          <input
            placeholder="Search nama / SKU"
            value={filters.q}
            onChange={(e) => setFilters((p) => ({ ...p, q: e.target.value }))}
          />
          <select value={filters.status} onChange={(e) => setFilters((p) => ({ ...p, status: e.target.value }))}>
            <option value="">Semua status</option>
            <option value="active">active</option>
            <option value="inactive">inactive</option>
          </select>
          <button type="submit">Cari</button>
          <button
            type="button"
            className="secondary"
            onClick={async () => {
              const reset = { q: '', status: '', page: 1, limit: 10 }
              setFilters(reset)
              await loadProducts(reset)
            }}
          >
            Reset
          </button>
        </form>

        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>SKU</th>
                <th>Nama</th>
                <th>Harga</th>
                <th>Status</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              {!loading && items.length === 0 && (
                <tr>
                  <td colSpan={6}>Belum ada data.</td>
                </tr>
              )}
              {items.map((item) => (
                <tr key={item.id}>
                  <td>{item.id}</td>
                  <td>{item.sku}</td>
                  <td>{item.name}</td>
                  <td>{Number(item.price).toLocaleString('id-ID')}</td>
                  <td>{item.status}</td>
                  <td className="row-actions">
                    <button className="secondary" onClick={() => handleEdit(item)}>Edit</button>
                    <button className="danger" onClick={() => handleDelete(item.id)}>Hapus</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className="pagination">
          <button
            className="secondary"
            disabled={meta.page <= 1}
            onClick={() => goToPage(meta.page - 1)}
          >
            Prev
          </button>
          <span>
            Page {meta.page} / {Math.max(meta.total_pages || 1, 1)} (Total: {meta.total || 0})
          </span>
          <button
            className="secondary"
            disabled={meta.page >= (meta.total_pages || 1)}
            onClick={() => goToPage(meta.page + 1)}
          >
            Next
          </button>
        </div>
      </section>
    </main>
  )
}

export default App
