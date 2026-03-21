import { useEffect, useMemo, useState } from 'react';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

async function api(path, options = {}) {
  const response = await fetch(`${API_URL}${path}`, {
    headers: { 'Content-Type': 'application/json', ...(options.headers || {}) },
    ...options,
  });

  const text = await response.text();
  const contentType = response.headers.get('content-type') || '';
  let data = null;

  if (text && contentType.includes('application/json')) {
    try {
      data = JSON.parse(text);
    } catch {
      data = null;
    }
  }

  if (!response.ok) {
    throw new Error((data && data.error) || text || `HTTP ${response.status}`);
  }

  return data;
}

function formatDate(value) {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '-';
  return date.toLocaleString();
}

export default function App() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const [stats, setStats] = useState(null);
  const [assets, setAssets] = useState([]);
  const [pagination, setPagination] = useState({ page: 1, limit: 20, total: 0, total_pages: 0 });

  const [filters, setFilters] = useState({ page: 1, limit: 20, type: '', status: '', q: '' });

  const [createForm, setCreateForm] = useState({ name: '', type: 'domain' });

  const [scanForm, setScanForm] = useState({ assetId: '', scanType: 'dns' });
  const [scanJobs, setScanJobs] = useState([]);
  const [selectedJobId, setSelectedJobId] = useState('');
  const [selectedJobStatus, setSelectedJobStatus] = useState(null);
  const [selectedJobResults, setSelectedJobResults] = useState(null);

  const scanTypeOptions = useMemo(
    () => ['dns', 'whois', 'subdomain', 'cert_trans', 'asn', 'all', 'ip', 'port', 'ssl', 'tech'],
    []
  );

  async function loadStats() {
    const data = await api('/assets/stats');
    setStats(data);
  }

  async function loadAssets(nextFilters = filters) {
    const params = new URLSearchParams();
    params.set('page', String(nextFilters.page));
    params.set('limit', String(nextFilters.limit));
    if (nextFilters.type) params.set('type', nextFilters.type);
    if (nextFilters.status) params.set('status', nextFilters.status);

    const data = await api(`/assets?${params.toString()}`);
    let list = data.data || [];

    if (nextFilters.q) {
      const searchData = await api(`/assets/search?q=${encodeURIComponent(nextFilters.q)}`);
      const allowed = new Set((searchData || []).map((asset) => asset.id));
      list = list.filter((asset) => allowed.has(asset.id));
    }

    setAssets(list);
    setPagination(data.pagination || { page: 1, limit: 20, total: 0, total_pages: 0 });
  }

  async function loadInitial() {
    setLoading(true);
    setError('');
    try {
      await Promise.all([loadStats(), loadAssets()]);
    } catch (err) {
      setError(err.message || 'Failed to load data');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadInitial();
  }, []);

  async function handleCreateAsset(event) {
    event.preventDefault();
    setError('');

    try {
      await api('/assets', {
        method: 'POST',
        body: JSON.stringify(createForm),
      });
      setCreateForm({ name: '', type: createForm.type });
      await Promise.all([loadStats(), loadAssets()]);
    } catch (err) {
      setError(err.message || 'Create asset failed');
    }
  }

  async function handleDeleteAsset(id) {
    setError('');
    try {
      await api(`/assets/${id}`, { method: 'DELETE' });
      await Promise.all([loadStats(), loadAssets()]);
    } catch (err) {
      setError(err.message || 'Delete asset failed');
    }
  }

  async function handleApplyFilters(event) {
    event.preventDefault();
    setError('');
    try {
      await loadAssets(filters);
    } catch (err) {
      setError(err.message || 'Filter failed');
    }
  }

  async function handlePageChange(nextPage) {
    const next = { ...filters, page: nextPage };
    setFilters(next);
    setError('');
    try {
      await loadAssets(next);
    } catch (err) {
      setError(err.message || 'Load page failed');
    }
  }

  async function handleStartScan(event) {
    event.preventDefault();
    if (!scanForm.assetId) return;

    setError('');
    try {
      const job = await api(`/assets/${scanForm.assetId}/scan`, {
        method: 'POST',
        body: JSON.stringify({ scan_type: scanForm.scanType }),
      });

      setScanJobs((prev) => [job, ...prev].slice(0, 10));
      setSelectedJobId(job.id);
      setSelectedJobStatus(job);
      setSelectedJobResults(null);
    } catch (err) {
      setError(err.message || 'Start scan failed');
    }
  }

  async function handleCheckJob() {
    if (!selectedJobId) return;
    setError('');

    try {
      const status = await api(`/scan-jobs/${selectedJobId}`);
      setSelectedJobStatus(status);
    } catch (err) {
      setError(err.message || 'Get job status failed');
    }
  }

  async function handleGetResults() {
    if (!selectedJobId) return;
    setError('');

    try {
      const results = await api(`/scan-jobs/${selectedJobId}/results`);
      setSelectedJobResults(results);
    } catch (err) {
      setError(err.message || 'Get results failed');
    }
  }

  return (
    <div className="page-shell">
      <div className="background-glow glow-a" />
      <div className="background-glow glow-b" />

      <header className="hero">
        <p className="hero-badge">EASM Frontend - Vite React</p>
        <h1>Asset Intelligence Command Center</h1>
        <p>
          Manage assets, launch scans, and inspect real-time reconnaissance results from your backend API.
        </p>
      </header>

      {error ? <div className="alert">{error}</div> : null}

      <section className="panel stats-panel">
        <div className="panel-header">
          <h2>Assets Statistics</h2>
          <button className="ghost-btn" onClick={loadInitial}>Refresh</button>
        </div>

        {loading ? (
          <p>Loading data...</p>
        ) : (
          <div className="stats-grid">
            <article className="stat-card">
              <p>Total Assets</p>
              <strong>{stats?.total ?? 0}</strong>
            </article>
            <article className="stat-card">
              <p>By Type</p>
              <strong>
                {Object.entries(stats?.by_type || {})
                  .map(([k, v]) => `${k}:${v}`)
                  .join(' | ') || '-'}
              </strong>
            </article>
            <article className="stat-card">
              <p>By Status</p>
              <strong>
                {Object.entries(stats?.by_status || {})
                  .map(([k, v]) => `${k}:${v}`)
                  .join(' | ') || '-'}
              </strong>
            </article>
          </div>
        )}
      </section>

      <section className="workspace-grid">
        <div className="panel">
          <h2>Create Asset</h2>
          <form className="stack-form" onSubmit={handleCreateAsset}>
            <label>
              Name
              <input
                required
                value={createForm.name}
                onChange={(event) => setCreateForm((prev) => ({ ...prev, name: event.target.value }))}
                placeholder="example.com or 127.0.0.1"
              />
            </label>

            <label>
              Type
              <select
                value={createForm.type}
                onChange={(event) => setCreateForm((prev) => ({ ...prev, type: event.target.value }))}
              >
                <option value="domain">domain</option>
                <option value="ip">ip</option>
                <option value="service">service</option>
              </select>
            </label>

            <button type="submit">Create</button>
          </form>
        </div>

        <div className="panel">
          <h2>Scan Control</h2>
          <form className="stack-form" onSubmit={handleStartScan}>
            <label>
              Asset
              <select
                value={scanForm.assetId}
                onChange={(event) => setScanForm((prev) => ({ ...prev, assetId: event.target.value }))}
              >
                <option value="">Select asset...</option>
                {assets.map((asset) => (
                  <option key={asset.id} value={asset.id}>
                    {asset.name} ({asset.type})
                  </option>
                ))}
              </select>
            </label>

            <label>
              Scan Type
              <select
                value={scanForm.scanType}
                onChange={(event) => setScanForm((prev) => ({ ...prev, scanType: event.target.value }))}
              >
                {scanTypeOptions.map((type) => (
                  <option key={type} value={type}>
                    {type}
                  </option>
                ))}
              </select>
            </label>

            <button type="submit">Start Scan</button>
          </form>

          <div className="scan-actions">
            <label>
              Job ID
              <input
                value={selectedJobId}
                onChange={(event) => setSelectedJobId(event.target.value)}
                placeholder="Paste scan job id"
              />
            </label>
            <div className="inline-actions">
              <button onClick={handleCheckJob}>Check Status</button>
              <button className="ghost-btn" onClick={handleGetResults}>Get Results</button>
            </div>
          </div>

          {selectedJobStatus ? (
            <pre className="json-box">{JSON.stringify(selectedJobStatus, null, 2)}</pre>
          ) : null}
        </div>
      </section>

      <section className="panel">
        <div className="panel-header">
          <h2>Assets Explorer</h2>
          <p>
            Page {pagination.page} / {pagination.total_pages || 1} - Total {pagination.total}
          </p>
        </div>

        <form className="filters" onSubmit={handleApplyFilters}>
          <input
            placeholder="Search by name (q)"
            value={filters.q}
            onChange={(event) => setFilters((prev) => ({ ...prev, q: event.target.value, page: 1 }))}
          />
          <select
            value={filters.type}
            onChange={(event) => setFilters((prev) => ({ ...prev, type: event.target.value, page: 1 }))}
          >
            <option value="">All types</option>
            <option value="domain">domain</option>
            <option value="ip">ip</option>
            <option value="service">service</option>
          </select>
          <select
            value={filters.status}
            onChange={(event) => setFilters((prev) => ({ ...prev, status: event.target.value, page: 1 }))}
          >
            <option value="">All status</option>
            <option value="active">active</option>
            <option value="inactive">inactive</option>
          </select>
          <select
            value={filters.limit}
            onChange={(event) => setFilters((prev) => ({ ...prev, limit: Number(event.target.value), page: 1 }))}
          >
            <option value={10}>10</option>
            <option value={20}>20</option>
            <option value={50}>50</option>
            <option value={100}>100</option>
          </select>
          <button type="submit">Apply</button>
        </form>

        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Type</th>
                <th>Status</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {assets.length === 0 ? (
                <tr>
                  <td colSpan={5}>No assets in current page/filter.</td>
                </tr>
              ) : (
                assets.map((asset) => (
                  <tr key={asset.id}>
                    <td>{asset.name}</td>
                    <td>{asset.type}</td>
                    <td>{asset.status}</td>
                    <td>{formatDate(asset.created_at)}</td>
                    <td>
                      <button className="danger-btn" onClick={() => handleDeleteAsset(asset.id)}>
                        Delete
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        <div className="pagination-row">
          <button
            disabled={filters.page <= 1}
            onClick={() => handlePageChange(filters.page - 1)}
          >
            Prev
          </button>
          <span>Current page: {filters.page}</span>
          <button
            disabled={pagination.total_pages > 0 && filters.page >= pagination.total_pages}
            onClick={() => handlePageChange(filters.page + 1)}
          >
            Next
          </button>
        </div>
      </section>

      {selectedJobResults ? (
        <section className="panel">
          <h2>Scan Results</h2>
          <pre className="json-box">{JSON.stringify(selectedJobResults, null, 2)}</pre>
        </section>
      ) : null}

      {scanJobs.length > 0 ? (
        <section className="panel">
          <h2>Latest Started Jobs</h2>
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Job ID</th>
                  <th>Asset ID</th>
                  <th>Type</th>
                  <th>Status</th>
                  <th>Created</th>
                </tr>
              </thead>
              <tbody>
                {scanJobs.map((job) => (
                  <tr key={job.id}>
                    <td>{job.id}</td>
                    <td>{job.asset_id}</td>
                    <td>{job.scan_type}</td>
                    <td>{job.status}</td>
                    <td>{formatDate(job.created_at)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>
      ) : null}
    </div>
  );
}
