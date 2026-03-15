import { useState } from 'react'
import { ChevronDown, ChevronRight, RotateCcw, FileText } from 'lucide-react'
import type { PromptSectionsConfig } from '../../types'

interface PromptSectionsEditorProps {
  config: PromptSectionsConfig | undefined
  onChange: (config: PromptSectionsConfig) => void
  disabled?: boolean
  language: string
}

// Default prompt sections (same as backend defaults)
const defaultSections: PromptSectionsConfig = {
  role_definition: `# You are a professional crypto trading AI

You focus on technical analysis and risk management, making rational trading decisions based on market data.
Your goal is to capture high-probability trading opportunities while controlling risk.`,

  trading_frequency: `# ⏱️ Trading Frequency Awareness

- Skilled traders: 2-4 trades/day ≈ 0.1-0.2 trades/hour
- >2 trades/hour = overtrading
- Minimum hold time per position ≥ 30-60 minutes
If you find yourself trading every cycle → standards too low; if closing positions within <30 minutes → too impatient.`,

  entry_standards: `# 🎯 Position Opening Criteria (Strict)

Only open positions when multiple signals converge:
- Clear trend direction (EMA alignment, price position)
- Momentum confirmation (MACD, RSI in agreement)
- Moderate volatility (ATR within reasonable range)
- Volume-price alignment (volume supports direction)

Avoid: single indicator, conflicting signals, ranging/choppy markets, re-entering immediately after closing.`,

  decision_process: `# 📋 Decision Process

1. Check positions → should we take profit / stop loss
2. Scan candidate coins + multiple timeframes → are there strong signals
3. Evaluate risk-reward ratio → does it meet minimum requirements
4. Write chain-of-thought first, then output structured JSON`,
}

export function PromptSectionsEditor({
  config,
  onChange,
  disabled,
  language,
}: PromptSectionsEditorProps) {
  const [expandedSections, setExpandedSections] = useState<Record<string, boolean>>({
    role_definition: false,
    trading_frequency: false,
    entry_standards: false,
    decision_process: false,
  })

  const t = (key: string) => {
    const translations: Record<string, Record<string, string>> = {
      promptSections: { zh: 'System Prompt 自定义', en: 'System Prompt Customization' },
      promptSectionsDesc: { zh: '自定义 AI 行为和决策逻辑（输出格式和风控规则不可修改）', en: 'Customize AI behavior and decision logic (output format and risk rules are fixed)' },
      roleDefinition: { zh: '角色定义', en: 'Role Definition' },
      roleDefinitionDesc: { zh: '定义 AI 的身份和核心目标', en: 'Define AI identity and core objectives' },
      tradingFrequency: { zh: '交易频率', en: 'Trading Frequency' },
      tradingFrequencyDesc: { zh: '设定交易频率预期和过度交易警告', en: 'Set trading frequency expectations and overtrading warnings' },
      entryStandards: { zh: '开仓标准', en: 'Entry Standards' },
      entryStandardsDesc: { zh: '定义开仓信号条件和避免事项', en: 'Define entry signal conditions and avoidances' },
      decisionProcess: { zh: '决策流程', en: 'Decision Process' },
      decisionProcessDesc: { zh: '设定决策步骤和思考流程', en: 'Set decision steps and thinking process' },
      resetToDefault: { zh: '重置为默认', en: 'Reset to Default' },
      chars: { zh: '字符', en: 'chars' },
    }
    return translations[key]?.[language] || key
  }

  const sections = [
    { key: 'role_definition', label: t('roleDefinition'), desc: t('roleDefinitionDesc') },
    { key: 'trading_frequency', label: t('tradingFrequency'), desc: t('tradingFrequencyDesc') },
    { key: 'entry_standards', label: t('entryStandards'), desc: t('entryStandardsDesc') },
    { key: 'decision_process', label: t('decisionProcess'), desc: t('decisionProcessDesc') },
  ]

  const currentConfig = config || {}

  const updateSection = (key: keyof PromptSectionsConfig, value: string) => {
    if (!disabled) {
      onChange({ ...currentConfig, [key]: value })
    }
  }

  const resetSection = (key: keyof PromptSectionsConfig) => {
    if (!disabled) {
      onChange({ ...currentConfig, [key]: defaultSections[key] })
    }
  }

  const toggleSection = (key: string) => {
    setExpandedSections((prev) => ({ ...prev, [key]: !prev[key] }))
  }

  const getValue = (key: keyof PromptSectionsConfig): string => {
    return currentConfig[key] || defaultSections[key] || ''
  }

  return (
    <div className="space-y-4">
      <div className="flex items-start gap-2 mb-4">
        <FileText className="w-5 h-5 mt-0.5" style={{ color: '#a855f7' }} />
        <div>
          <h3 className="font-medium" style={{ color: '#EAECEF' }}>
            {t('promptSections')}
          </h3>
          <p className="text-xs mt-1" style={{ color: '#848E9C' }}>
            {t('promptSectionsDesc')}
          </p>
        </div>
      </div>

      <div className="space-y-2">
        {sections.map(({ key, label, desc }) => {
          const sectionKey = key as keyof PromptSectionsConfig
          const isExpanded = expandedSections[key]
          const value = getValue(sectionKey)
          const isModified = currentConfig[sectionKey] !== undefined && currentConfig[sectionKey] !== defaultSections[sectionKey]

          return (
            <div
              key={key}
              className="rounded-lg overflow-hidden"
              style={{ background: '#0B0E11', border: '1px solid #2B3139' }}
            >
              <button
                onClick={() => toggleSection(key)}
                className="w-full flex items-center justify-between px-3 py-2.5 hover:bg-white/5 transition-colors text-left"
              >
                <div className="flex items-center gap-2">
                  {isExpanded ? (
                    <ChevronDown className="w-4 h-4" style={{ color: '#848E9C' }} />
                  ) : (
                    <ChevronRight className="w-4 h-4" style={{ color: '#848E9C' }} />
                  )}
                  <span className="text-sm font-medium" style={{ color: '#EAECEF' }}>
                    {label}
                  </span>
                  {isModified && (
                    <span
                      className="px-1.5 py-0.5 text-[10px] rounded"
                      style={{ background: 'rgba(168, 85, 247, 0.15)', color: '#a855f7' }}
                    >
                      {language === 'zh' ? '已修改' : 'Modified'}
                    </span>
                  )}
                </div>
                <span className="text-[10px]" style={{ color: '#848E9C' }}>
                  {value.length} {t('chars')}
                </span>
              </button>

              {isExpanded && (
                <div className="px-3 pb-3">
                  <p className="text-xs mb-2" style={{ color: '#848E9C' }}>
                    {desc}
                  </p>
                  <textarea
                    value={value}
                    onChange={(e) => updateSection(sectionKey, e.target.value)}
                    disabled={disabled}
                    rows={6}
                    className="w-full px-3 py-2 rounded-lg resize-y font-mono text-xs"
                    style={{
                      background: '#1E2329',
                      border: '1px solid #2B3139',
                      color: '#EAECEF',
                      minHeight: '120px',
                    }}
                  />
                  <div className="flex justify-end mt-2">
                    <button
                      onClick={() => resetSection(sectionKey)}
                      disabled={disabled || !isModified}
                      className="flex items-center gap-1 px-2 py-1 rounded text-xs transition-colors hover:bg-white/5 disabled:opacity-30"
                      style={{ color: '#848E9C' }}
                    >
                      <RotateCcw className="w-3 h-3" />
                      {t('resetToDefault')}
                    </button>
                  </div>
                </div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
