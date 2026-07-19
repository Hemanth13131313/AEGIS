import React from 'react';

export function ProvidersPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">AI Provider Adapters</h1>
      
      <div className="bg-blue-50 text-blue-800 p-4 rounded-md mb-6 border border-blue-200">
        <p>Provider credentials are stored in Vault — not in this UI.</p>
      </div>

      <div className="flex justify-end mb-6">
        <button className="bg-blue-600 text-white px-4 py-2 rounded shadow hover:bg-blue-700 transition">
          Add Custom Provider
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
        {[
          { name: 'OpenAI', status: 'Configured', endpoint: 'https://api.openai.com/v1' },
          { name: 'Anthropic', status: 'Configured', endpoint: 'https://api.anthropic.com/v1' },
          { name: 'Google Gemini', status: 'Not Configured', endpoint: 'https://generativelanguage.googleapis.com' },
          { name: 'AWS Bedrock', status: 'Not Configured', endpoint: 'https://bedrock-runtime.us-east-1.amazonaws.com' },
        ].map((p) => (
          <div key={p.name} className="border rounded shadow-sm bg-white p-4 flex flex-col">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-lg font-semibold">{p.name}</h2>
              <span className={`px-2 py-1 text-xs font-semibold rounded-full ${
                p.status === 'Configured' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
              }`}>
                {p.status}
              </span>
            </div>
            
            <div className="mb-4 flex-grow">
              <label className="block text-sm text-gray-500 mb-1">Endpoint</label>
              <input 
                type="text" 
                readOnly 
                value={p.endpoint} 
                className="w-full bg-gray-50 border border-gray-200 rounded p-2 text-sm text-gray-700" 
              />
            </div>

            <button className="w-full border border-gray-300 text-gray-700 rounded py-2 hover:bg-gray-50 transition">
              Configure
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
