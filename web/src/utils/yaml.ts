import { parse, stringify } from 'yaml'

/**
 * Parse a YAML document into a plain config object.
 *
 * Backed by the `yaml` library so flow-style values (`[200, 401, 404]`,
 * `{ a: 1 }`), arbitrarily nested structures, and quoting all round-trip
 * correctly. Throws on invalid YAML — callers are expected to surface the
 * error to the user.
 *
 * Route/middleware configs are always top-level maps, so non-map documents
 * (a bare scalar or sequence) collapse to an empty object rather than
 * producing a surprising shape downstream.
 */
export function parseYaml(text: string): Record<string, unknown> {
  const result = parse(text)
  if (result && typeof result === 'object' && !Array.isArray(result)) {
    return result as Record<string, unknown>
  }
  return {}
}

/**
 * Serialize a config object back to YAML for the editor.
 *
 * `lineWidth: 0` disables line folding so long scalars (URLs, regexes) stay
 * on a single line.
 */
export function toYaml(obj: unknown): string {
  if (obj === undefined || obj === null) return ''
  return stringify(obj, { lineWidth: 0 })
}
