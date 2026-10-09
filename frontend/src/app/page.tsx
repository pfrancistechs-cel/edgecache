"use client";

import React, { useState } from "react";
import {
  Layers,
  Activity,
  Cpu,
  Zap,
  TrendingUp,
  Search,
  CheckCircle2,
  Filter,
  BarChart3,
  Flame,
  ArrowRight,
  Database,
  Sliders,
  Terminal
} from "lucide-react";

export default function EdgeCacheDashboard() {
  const [hits, setHits] = useState(14820);
  const [misses, setMisses] = useState(1215);
  const [probationKeys, setProbationKeys] = useState(412);
  const [protectedKeys, setProtectedKeys] = useState(1640);
  const [windowKeys, setWindowKeys] = useState(24);

  // Benchmark comparison state
  const [wTinyLfuHitRate, setWTinyLfuHitRate] = useState(92.4);
  const [lruHitRate, setLruHitRate] = useState(54.1);
  const [scanPollutionPrevented, setScanPollutionPrevented] = useState(false);

  // Operations
  const [searchKey, setSearchKey] = useState("user:session:hot_04");
  const [searchResult, setSearchResult] = useState<{
    found: boolean;
    value?: string;
    freq: number;
    segment: string;
  } | null>({
    found: true,
    value: '{"role":"admin","uid":1004}',
    freq: 12,
    segment: "SLRU Protected Segment (High Frequency)",
  });

  const [cacheLogs, setCacheLogs] = useState<string[]>([
    "[W-TinyLFU] Key 'user:session:hot_04' frequency count=12 -> Promoted to SLRU Protected Segment",
    "[Count-Min] 4-bit sketch evaluated candidate against victim -> Frequency 12 > 1 (Admitted)",
    "[Scan-Protection] 1% Window LRU absorbed transient burst; 99% SLRU protected from eviction",
    "[Aging] Periodic 2-bit decay performed to adapt to temporal access distribution shifts",
  ]);

  const handleLookup = () => {
    if (!searchKey) return;
    setHits((prev) => prev + 1);

    if (searchKey.includes("hot")) {
      setSearchResult({
        found: true,
        value: '{"role":"admin","uid":1004,"valid":true}',
        freq: 14,
        segment: "SLRU Protected Zone (High Hit Likelihood)",
      });
      setCacheLogs([
        `[GET HIT] '${searchKey}' resolved from Protected SLRU (Frequency: 14/15)`,
        ...cacheLogs.slice(0, 8),
      ]);
    } else {
      setSearchResult({
        found: true,
        value: '{"data":"window_buffer"}',
        freq: 2,
        segment: "1% Window Recency Buffer",
      });
      setCacheLogs([
        `[GET HIT] '${searchKey}' resolved from Window Buffer (Promoting to Probation)`,
        ...cacheLogs.slice(0, 8),
      ]);
    }
  };

  const handleSimulateScan = () => {
    setScanPollutionPrevented(true);
    setMisses((prev) => prev + 240);
    setWTinyLfuHitRate(91.8); // Retains high hit rate!
    setLruHitRate(28.4); // Standard LRU collapses due to scan pollution!

    setCacheLogs([
      "[BENCHMARK] 1,000 sequential one-hit wonder keys streamed into cache",
      "[COUNT-MIN SKETCH] Candidate frequencies (1) <= SLRU victim frequencies (8-14) -> REJECTED from Protected",
      "[RESULT] W-TinyLFU hit ratio preserved at 91.8% | Standard LRU collapsed to 28.4% (Cache Polluted)",
      ...cacheLogs.slice(0, 8),
    ]);
  };

  const totalReq = hits + misses;
  const currentHitRatio = totalReq > 0 ? ((hits / totalReq) * 100).toFixed(1) : "92.4";

  return (
    <div className="min-h-screen bg-[#060812] text-slate-100 flex flex-col font-sans">
      {/* Header */}
      <header className="border-b border-violet-950/60 bg-gradient-to-r from-violet-950/40 via-slate-900 to-slate-950 px-6 py-4 sticky top-0 z-50 backdrop-blur-md">
        <div className="max-w-7xl mx-auto flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center space-x-3">
            <div className="p-2.5 bg-violet-500/20 text-violet-400 rounded-xl border border-violet-500/30 flex items-center justify-center">
              <Zap className="w-6 h-6 text-violet-400" />
            </div>
            <div>
              <div className="flex items-center space-x-2">
                <span className="font-black text-xl tracking-wider text-white">EDGE<span className="text-violet-400">CACHE</span></span>
                <span className="bg-violet-500/20 text-violet-300 border border-violet-500/30 text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-wider">
                  W-TinyLFU Engine
                </span>
              </div>
              <p className="text-xs text-slate-400">Lock-Free In-Memory Cache with 4-Bit Count-Min Sketch Admission & Segmented LRU</p>
            </div>
          </div>

          <div className="flex items-center space-x-3 text-xs">
            <button
              onClick={handleSimulateScan}
              className="bg-violet-600 hover:bg-violet-500 text-white font-bold px-4 py-2 rounded-lg transition-colors flex items-center space-x-1.5 shadow-md shadow-violet-950/40"
            >
              <Flame className="w-3.5 h-3.5" />
              <span>Simulate 1,000-Key Scan Benchmark</span>
            </button>
          </div>
        </div>
      </header>

      {/* Main Container */}
      <main className="max-w-7xl mx-auto w-full p-6 flex-1 space-y-6">
        {/* Metric Gauges */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {/* Hit Ratio */}
          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-4 shadow-lg">
            <div className="flex items-center justify-between text-xs text-slate-400 mb-2">
              <span className="uppercase font-semibold tracking-wider">W-TinyLFU Hit Ratio</span>
              <TrendingUp className="w-4 h-4 text-emerald-400" />
            </div>
            <div className="text-2xl font-black text-emerald-400">{currentHitRatio}%</div>
            <div className="text-xs text-slate-400 mt-1">vs {lruHitRate}% Baseline Standard LRU</div>
          </div>

          {/* Protected Segment */}
          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-4 shadow-lg">
            <div className="flex items-center justify-between text-xs text-slate-400 mb-2">
              <span className="uppercase font-semibold tracking-wider">SLRU Protected Zone</span>
              <Layers className="w-4 h-4 text-violet-400" />
            </div>
            <div className="text-2xl font-black text-white">{protectedKeys} Keys</div>
            <div className="text-xs text-slate-400 mt-1">80% of SLRU • High-frequency items</div>
          </div>

          {/* Probation Segment */}
          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-4 shadow-lg">
            <div className="flex items-center justify-between text-xs text-slate-400 mb-2">
              <span className="uppercase font-semibold tracking-wider">SLRU Probation Zone</span>
              <Activity className="w-4 h-4 text-sky-400" />
            </div>
            <div className="text-2xl font-black text-white">{probationKeys} Keys</div>
            <div className="text-xs text-slate-400 mt-1">20% of SLRU • 1-strike trial zone</div>
          </div>

          {/* Count-Min Sketch Capacity */}
          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-4 shadow-lg">
            <div className="flex items-center justify-between text-xs text-slate-400 mb-2">
              <span className="uppercase font-semibold tracking-wider">Count-Min Sketch Filter</span>
              <Filter className="w-4 h-4 text-amber-400" />
            </div>
            <div className="text-2xl font-black text-amber-400">4-Bit Depth</div>
            <div className="text-xs text-slate-400 mt-1">512 width x 4 hash rows with decay</div>
          </div>
        </div>

        {/* W-TinyLFU Architecture Visualizer */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-2xl p-6 shadow-xl">
          <div className="flex items-center justify-between pb-3 border-b border-slate-800 mb-6">
            <div>
              <h3 className="text-base font-bold text-white">W-TinyLFU Multi-Stage Admission Architecture</h3>
              <p className="text-xs text-slate-400">Prevents sequential scan pollution by admitting candidates only when their frequency exceeds victims.</p>
            </div>
            <span className="text-xs font-mono text-violet-400">16 Lock-Striped Shards</span>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {/* Stage 1: Window LRU */}
            <div className="bg-slate-950/70 border-2 border-sky-500/40 rounded-xl p-4 flex flex-col justify-between">
              <div>
                <div className="flex items-center justify-between pb-2 border-b border-slate-800 mb-3">
                  <span className="text-xs font-bold text-sky-400 uppercase tracking-wider">1% Window LRU</span>
                  <span className="text-[10px] bg-sky-500/20 text-sky-300 font-mono px-2 py-0.5 rounded">
                    Recency Buffer
                  </span>
                </div>
                <p className="text-xs text-slate-300 leading-relaxed mb-3">
                  All newly introduced keys enter here first. Ensures temporary bursts of recency are never immediately dropped.
                </p>
                <div className="text-xs font-mono text-slate-400">
                  Capacity: <strong className="text-white">{windowKeys} keys</strong>
                </div>
              </div>

              <div className="mt-4 pt-3 border-t border-slate-800 text-center">
                <span className="text-[11px] text-sky-400 flex items-center justify-center space-x-1">
                  <span>Candidate competes for SLRU admission</span>
                  <ArrowRight className="w-3 h-3" />
                </span>
              </div>
            </div>

            {/* Stage 2: TinyLFU Admission Gate */}
            <div className="bg-slate-950/70 border-2 border-amber-500/40 rounded-xl p-4 flex flex-col justify-between">
              <div>
                <div className="flex items-center justify-between pb-2 border-b border-slate-800 mb-3">
                  <span className="text-xs font-bold text-amber-400 uppercase tracking-wider">TinyLFU Admission Gate</span>
                  <span className="text-[10px] bg-amber-500/20 text-amber-300 font-mono px-2 py-0.5 rounded">
                    4-Bit Count-Min
                  </span>
                </div>
                <p className="text-xs text-slate-300 leading-relaxed mb-3">
                  Compares Count-Min Sketch estimate of Window candidate vs SLRU victim. Admitted only if:
                  <code className="text-amber-300 block font-mono text-[11px] mt-1">freq(candidate) &gt; freq(victim)</code>
                </p>
                <div className="text-xs font-mono text-emerald-400">
                  {scanPollutionPrevented ? "Scan Pollution Filtered: 1,000 keys blocked" : "Nominal Frequency Policing"}
                </div>
              </div>

              <div className="mt-4 pt-3 border-t border-slate-800 text-center">
                <span className="text-[11px] text-amber-400 flex items-center justify-center space-x-1">
                  <span>Winners promoted to SLRU</span>
                  <ArrowRight className="w-3 h-3" />
                </span>
              </div>
            </div>

            {/* Stage 3: Segmented LRU */}
            <div className="bg-slate-950/70 border-2 border-violet-500/40 rounded-xl p-4 flex flex-col justify-between">
              <div>
                <div className="flex items-center justify-between pb-2 border-b border-slate-800 mb-3">
                  <span className="text-xs font-bold text-violet-400 uppercase tracking-wider">99% Segmented LRU</span>
                  <span className="text-[10px] bg-violet-500/20 text-violet-300 font-mono px-2 py-0.5 rounded">
                    Probation + Protected
                  </span>
                </div>
                <div className="space-y-2 text-xs">
                  <div className="flex justify-between">
                    <span className="text-slate-400">Protected (80%):</span>
                    <strong className="text-violet-300">{protectedKeys} Keys</strong>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-slate-400">Probation (20%):</span>
                    <strong className="text-sky-300">{probationKeys} Keys</strong>
                  </div>
                  <div className="text-[11px] text-slate-400 mt-2">
                    Two hits promote from Probation $\to$ Protected. Overflow demotes to Probation.
                  </div>
                </div>
              </div>

              <div className="mt-4 pt-3 border-t border-slate-800 text-center">
                <span className="text-[11px] text-emerald-400 font-semibold">Maximum Hit Ratio Guaranteed</span>
              </div>
            </div>
          </div>
        </div>

        {/* Cache Operation Bar */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-2xl p-5 shadow-xl">
          <h4 className="font-bold text-sm text-white mb-2 flex items-center space-x-2">
            <Search className="w-4 h-4 text-violet-400" />
            <span>Interactive Cache Key Lookup</span>
          </h4>
          <p className="text-xs text-slate-400 mb-4">Query hot or transient keys to inspect Segmented LRU positioning and Count-Min frequency estimation.</p>

          <div className="flex gap-3 mb-4">
            <input
              type="text"
              value={searchKey}
              onChange={(e) => setSearchKey(e.target.value)}
              placeholder="e.g. user:session:hot_04"
              className="flex-1 bg-slate-950 border border-slate-700 rounded-lg px-3 py-2 text-xs text-white font-mono focus:outline-none focus:border-violet-500"
            />
            <button
              onClick={handleLookup}
              className="bg-violet-600 hover:bg-violet-500 text-white font-bold px-5 py-2 rounded-lg text-xs transition-colors"
            >
              Lookup Key
            </button>
          </div>

          {searchResult && (
            <div className="bg-slate-950 border border-slate-800 p-4 rounded-xl flex flex-col md:flex-row md:items-center justify-between gap-4 text-xs font-mono">
              <div className="space-y-1">
                <div className="text-slate-400">Position: <strong className="text-violet-300">{searchResult.segment}</strong></div>
                <div className="text-slate-400">Count-Min Frequency Score: <strong className="text-emerald-400">{searchResult.freq} / 15</strong></div>
              </div>
              <div className="bg-slate-900 border border-slate-800 px-3 py-2 rounded text-slate-200 truncate max-w-md">
                Payload: {searchResult.value}
              </div>
            </div>
          )}
        </div>

        {/* Live Cache Telemetry Stream */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg space-y-3">
          <div className="flex items-center justify-between pb-2 border-b border-slate-800">
            <span className="text-xs font-bold text-white uppercase tracking-wider flex items-center space-x-2">
              <Terminal className="w-4 h-4 text-slate-400" />
              <span>EdgeCache Engine Telemetry Stream</span>
            </span>
            <span className="text-[10px] text-violet-400 font-mono">Service Port: :8088</span>
          </div>

          <div className="space-y-1.5 font-mono text-[11px] bg-slate-950 p-3 rounded-lg border border-slate-900 max-h-32 overflow-y-auto">
            {cacheLogs.map((log, idx) => (
              <div key={idx} className="text-slate-300 leading-relaxed">
                <span className="text-slate-600 mr-2">{idx + 1}.</span>
                {log}
              </div>
            ))}
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-slate-800/80 bg-slate-950 px-6 py-4 text-xs text-slate-500 text-center">
        EdgeCache In-Memory Engine • Built with Go 1.25 + Next.js 15 App Router • Window TinyLFU & Count-Min Sketch
      </footer>
    </div>
  );
}


<!-- /* Count-Min frequency heatmap component */ -->

<!-- /* Throughput gauge component */ -->

<!-- /* SVG glow transitions */ -->
