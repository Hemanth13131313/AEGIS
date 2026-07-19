export const SEVERITY_LEVELS = ['critical', 'high', 'medium', 'low', 'safe'] as const;

export type Severity = typeof SEVERITY_LEVELS[number];

export const OWASP_CATEGORIES = [
  { id: 'LLM01', name: 'Prompt Injection' },
  { id: 'LLM02', name: 'Insecure Output Handling' },
  { id: 'LLM03', name: 'Training Data Poisoning' },
  { id: 'LLM04', name: 'Model Denial of Service' },
  { id: 'LLM05', name: 'Supply Chain Vulnerabilities' },
  { id: 'LLM06', name: 'Sensitive Information Disclosure' },
  { id: 'LLM07', name: 'Insecure Plugin Design' },
  { id: 'LLM08', name: 'Excessive Agency' },
  { id: 'LLM09', name: 'Overreliance' },
  { id: 'LLM10', name: 'Model Theft' }
] as const;

export const MITRE_ATLAS_TECHNIQUES = [
  { id: 'AML.T0001', name: 'Poison Training Data' },
  { id: 'AML.T0051', name: 'LLM Prompt Injection' },
  { id: 'AML.T0054', name: 'LLM Jailbreak' },
  { id: 'AML.T0055', name: 'Evade ML Model' },
  { id: 'AML.T0053', name: 'Command and Scripting Interpreter' }
] as const;
