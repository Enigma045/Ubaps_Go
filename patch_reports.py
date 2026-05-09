"""
Patch Dashboard.css with analytics strip styles
and Student_reports.js with fetchReportStats logic.
"""

# ── 1. CSS ───────────────────────────────────────────────────────────────────

CSS_APPEND = r"""

/* ═══════════════════════════════════════════════════
   REPORTS PAGE — compact layout overrides
   ═══════════════════════════════════════════════════ */

/* Remove padding from app so table reaches edges */
.app {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    padding: 0;
}

/* Allow inner sections their own side padding */
.tabs,
.filter-bar,
.pagination {
    padding-left: 1.5rem;
    padding-right: 1.5rem;
}

/* Table fills remaining height, no rounded corners */
.table-container {
    background: #fff;
    border-top: 1px solid #e2e8f0;
    overflow-y: auto;
    flex: 1;
    min-height: 0;
    width: 100%;
}

.log-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.8125rem;
    table-layout: fixed;
}

.log-table th {
    background: #f8fafc;
    color: #2563eb;
    font-weight: 700;
    text-align: left;
    padding: 6px 12px;
    border-bottom: 1px solid #e2e8f0;
    white-space: nowrap;
    position: sticky;
    top: 0;
    z-index: 10;
}

.log-table td {
    padding: 4px 12px;
    border-bottom: 1px solid #e2e8f0;
    color: #1e293b;
    vertical-align: middle;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.log-table tr:last-child td { border-bottom: none; }
.log-table tr:hover { background: #f1f5f9; }

/* ── Report Top Bar ───────────────────────────────── */
.report-topbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.55rem 1.5rem;
    background: #fff;
    border-bottom: 1px solid #e2e8f0;
    flex-shrink: 0;
    gap: 1rem;
}

.report-title-block h1 {
    font-size: 1rem;
    font-weight: 800;
    color: #1e293b;
    margin: 0;
    line-height: 1.2;
}

.report-title-block .subtitle {
    font-size: 0.68rem;
    color: #64748b;
    margin: 0;
}

.period-controls {
    display: flex;
    align-items: flex-end;
    gap: 0.6rem;
}

.period-item {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.period-item label {
    font-size: 0.58rem;
    font-weight: 700;
    color: #64748b;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}

.period-item select {
    padding: 0.28rem 0.55rem;
    border: 1px solid #e2e8f0;
    border-radius: 6px;
    font-size: 0.78rem;
    color: #1e293b;
    background: #f8fafc;
    cursor: pointer;
    outline: none;
    min-width: 110px;
}

.period-item select:focus {
    border-color: #2563eb;
    background: #fff;
    box-shadow: 0 0 0 3px rgba(37,99,235,0.1);
}

/* ── Analytics Strip ─────────────────────────────── */
.analytics-strip {
    display: flex;
    align-items: center;
    padding: 0 1.5rem;
    background: #f8fafc;
    border-bottom: 1px solid #e2e8f0;
    flex-shrink: 0;
    overflow-x: auto;
    min-height: 46px;
}

.astat {
    display: flex;
    flex-direction: column;
    justify-content: center;
    padding: 0.4rem 0.9rem;
    gap: 1px;
    white-space: nowrap;
    position: relative;
    cursor: default;
}

.astat-label {
    font-size: 0.58rem;
    font-weight: 700;
    color: #94a3b8;
    text-transform: uppercase;
    letter-spacing: 0.06em;
}

.astat-value {
    font-size: 1rem;
    font-weight: 800;
    color: #1e293b;
    letter-spacing: -0.01em;
    line-height: 1;
}

/* Coloured status dots */
.astat-dot {
    position: absolute;
    top: 0.5rem;
    left: 0.45rem;
    width: 6px;
    height: 6px;
    border-radius: 50%;
}

.astat.approved { padding-left: 1.35rem; }
.astat.approved .astat-dot  { background: #10b981; }
.astat.approved .astat-value { color: #065f46; }

.astat.pending  { padding-left: 1.35rem; }
.astat.pending  .astat-dot  { background: #f59e0b; }
.astat.pending  .astat-value { color: #92400e; }

.astat.rejected { padding-left: 1.35rem; }
.astat.rejected .astat-dot  { background: #ef4444; }
.astat.rejected .astat-value { color: #991b1b; }

.astat.total   .astat-value { color: #2563eb; }
.astat.value   .astat-value { color: #7c3aed; }
.astat.avgtime .astat-value { color: #0369a1; }

/* Vertical separators */
.astat-sep {
    width: 1px;
    height: 26px;
    background: #e2e8f0;
    flex-shrink: 0;
    margin: 0 0.25rem;
}

/* Faculty chips section */
.astat-faculty {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0 0.9rem;
    white-space: nowrap;
}

.astat-faculty > .astat-label {
    font-size: 0.58rem;
    font-weight: 700;
    color: #94a3b8;
    text-transform: uppercase;
    letter-spacing: 0.06em;
}

.faculty-chips {
    display: flex;
    gap: 0.35rem;
    flex-wrap: nowrap;
}

.faculty-chip {
    display: inline-flex;
    align-items: baseline;
    gap: 0.25rem;
    padding: 0.15rem 0.55rem;
    border-radius: 999px;
    background: #dbeafe;
    border: 1px solid #bfdbfe;
    font-size: 0.68rem;
    font-weight: 700;
    color: #1d4ed8;
    white-space: nowrap;
}

.faculty-chip .chip-count {
    font-weight: 900;
    font-size: 0.78rem;
    color: #1e40af;
}

.faculty-loading {
    font-size: 0.72rem;
    color: #94a3b8;
    font-style: italic;
}
"""

css_path = r'c:\Users\USER\Go\Ubaps\Pages\Html\student\Css\Dashboard.css'
with open(css_path, 'r', encoding='utf-8') as f:
    css = f.read()

# Remove previously appended reports-specific blocks to avoid duplication
# We'll detect from the first occurrence of our marker
marker = '/* ═══════════════════════════════════════════════════'
if marker in css:
    css = css[:css.index(marker)]

css = css.rstrip() + '\n' + CSS_APPEND

with open(css_path, 'w', encoding='utf-8') as f:
    f.write(css)

print('CSS patched OK')

# ── 2. JS ────────────────────────────────────────────────────────────────────

JS_EXTRA = r"""
async function fetchReportStats() {
    const year     = document.getElementById('period-year')?.value     || '';
    const semester = document.getElementById('period-semester')?.value || '';
    const month    = document.getElementById('period-month')?.value    || '';

    try {
        const res = await fetch(`/api/reports/stats?year=${year}&semester=${semester}&month=${month}`);
        if (!res.ok) throw new Error('Stats fetch failed');
        const s = await res.json();

        const set = (id, val) => { const el = document.getElementById(id); if (el) el.textContent = val; };

        set('stat-total',    s.total_applications ?? '—');
        set('stat-approved', s.approved           ?? '—');
        set('stat-pending',  s.pending            ?? '—');
        set('stat-rejected', s.rejected           ?? '—');
        set('stat-value',    s.total_value != null
            ? 'MWK ' + Number(s.total_value).toLocaleString(undefined, {maximumFractionDigits: 0})
            : '—');
        set('stat-avgtime',  s.avg_processing_time != null
            ? s.avg_processing_time + ' days'
            : '—');

        // Faculty chips
        const fc = document.getElementById('faculty-chips');
        if (fc) {
            fc.innerHTML = '';
            const breakdown = s.faculty_breakdown || [];
            if (breakdown.length === 0) {
                fc.innerHTML = '<span class="faculty-loading">No data</span>';
            } else {
                breakdown.forEach(item => {
                    const chip = document.createElement('span');
                    chip.className = 'faculty-chip';
                    chip.innerHTML = `${item.faculty} <span class="chip-count">${item.count}</span>`;
                    fc.appendChild(chip);
                });
            }
        }
    } catch (err) {
        console.warn('fetchReportStats error:', err);
    }
}
"""

JS_LISTENER = """
    // Period selectors trigger stats refresh
    ['period-year','period-semester','period-month'].forEach(id => {
        const el = document.getElementById(id);
        if (el) el.addEventListener('change', fetchReportStats);
    });

    // Initial stats load
    fetchReportStats();
"""

js_path = r'c:\Users\USER\Go\Ubaps\Pages\Html\student\Js\Student_reports.js'
with open(js_path, 'r', encoding='utf-8') as f:
    js = f.read()

# 1. Append fetchReportStats function (once)
if 'fetchReportStats' not in js:
    js = js.rstrip() + '\n' + JS_EXTRA

# 2. Inject listener + initial call inside DOMContentLoaded, before closing });
# Find the closing of DOMContentLoaded
target = '    filter.forEach(e => {\n        e.addEventListener("change", () => {\n            currentPage = 1; // Reset to page 1 on filter change\n            updatePayloadAndFetch();\n        });\n    });\n});'
replacement = target.replace('});\n});', '});\n' + JS_LISTENER + '\n});')

if JS_LISTENER.strip() not in js:
    if target in js:
        js = js.replace(target, replacement, 1)
    else:
        print('WARNING: DOMContentLoaded closing pattern not found — check JS manually')

with open(js_path, 'w', encoding='utf-8') as f:
    f.write(js)

print('JS patched OK')
print('All done.')
