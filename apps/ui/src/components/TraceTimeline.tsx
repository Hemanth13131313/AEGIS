import { TraceEvent } from '../hooks/useSessionTrace';
import { ArrowRight, Filter, Shield, Search, Layers, Server, Globe } from 'lucide-react';

interface Props {
  events: TraceEvent[];
  totalLatencyMs: number;
}

export function TraceTimeline({ events, totalLatencyMs }: Props) {
  const getIcon = (stage: string) => {
    switch (stage) {
      case 'request': return <ArrowRight className="w-5 h-5 text-blue-500" />;
      case 'pre_filter': return <Filter className="w-5 h-5 text-gray-500" />;
      case 'policy': return <Shield className="w-5 h-5 text-indigo-500" />;
      case 'scanner_a':
      case 'scanner_b': return <Search className="w-5 h-5 text-purple-500" />;
      case 'ensemble': return <Layers className="w-5 h-5 text-orange-500" />;
      case 'gateway': return <Server className="w-5 h-5 text-gray-700" />;
      case 'upstream': return <Globe className="w-5 h-5 text-green-500" />;
      default: return <div className="w-5 h-5 rounded-full bg-gray-300" />;
    }
  };

  const getVerdictColor = (verdict?: string) => {
    if (verdict === 'allow' || verdict === 'pass') return 'text-green-600 bg-green-50';
    if (verdict === 'block' || verdict === 'fail') return 'text-red-600 bg-red-50';
    if (verdict === 'tag') return 'text-orange-600 bg-orange-50';
    return 'text-gray-600 bg-gray-50';
  };

  const getLatencyColor = (ms: number) => {
    if (ms < 10) return 'text-green-600';
    if (ms < 100) return 'text-yellow-600';
    return 'text-red-600';
  };

  return (
    <div className="relative border-l-2 border-gray-200 ml-4 pl-6 space-y-6">
      {events.map((event) => (
        <div key={event.id} className="relative">
          <div className="absolute -left-9 top-1 bg-white p-1 rounded-full border border-gray-200">
            {getIcon(event.stage)}
          </div>
          
          <div className={`p-3 rounded border ${getVerdictColor(event.verdict)}`}>
            <div className="flex justify-between items-start mb-1">
              <div className="font-semibold">{event.label}</div>
              <div className="text-xs text-gray-500">{event.timestamp}</div>
            </div>
            <div className="text-sm">{event.details}</div>
            <div className={`text-xs mt-2 font-medium ${getLatencyColor(event.latency_ms)}`}>
              Latency: {event.latency_ms}ms
            </div>
          </div>
        </div>
      ))}
      <div className="mt-6 pt-4 border-t font-bold text-gray-700">
        Total pipeline: <span className={getLatencyColor(totalLatencyMs)}>{totalLatencyMs}ms</span>
      </div>
    </div>
  );
}
