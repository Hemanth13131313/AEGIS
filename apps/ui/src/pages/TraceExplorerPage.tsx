import { useState } from 'react';
import { useSessionTrace } from '../hooks/useSessionTrace';
import { TraceTimeline } from '../components/TraceTimeline';

export function TraceExplorerPage() {
  const [searchInput, setSearchInput] = useState('');
  const [sessionId, setSessionId] = useState<string | null>(null);

  const { data, isLoading } = useSessionTrace(sessionId);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSessionId(searchInput);
  };

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">Session Trace Explorer</h1>
      <form onSubmit={handleSearch} className="mb-8 flex gap-2">
        <input
          type="text"
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          placeholder="Enter Session ID (e.g., sess-001)"
          className="border p-2 rounded flex-1"
        />
        <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded">
          Search
        </button>
      </form>

      {isLoading && <p>Loading trace...</p>}
      
      {!isLoading && !data && sessionId && (
        <div className="p-4 bg-gray-100 rounded text-center">
          <p>Session not found. Try "sess-001" or "sess-002".</p>
        </div>
      )}

      {!isLoading && data && (
        <div className="flex flex-col gap-6">
          <div className="bg-gray-50 border p-4 rounded font-mono text-sm">
            <h3 className="font-bold text-gray-700 mb-2">Payload Preview</h3>
            <div className="text-gray-500 italic">[REDACTED — payload not stored]</div>
          </div>
          
          <div className="bg-white border rounded p-6 shadow-sm">
            <TraceTimeline events={data.events} totalLatencyMs={data.total_latency_ms} />
          </div>
        </div>
      )}
    </div>
  );
}
