export const gameCatalog = [
  ['word_hunter', 'Word Hunter', 'Conecta palabras con imágenes.', 'vocabulary'],
  ['sentence_builder', 'Sentence Builder', 'Dale orden a tus ideas.', 'grammar'],
  ['memory_match', 'Memory Match', 'Recuerda y relaciona significados.', 'vocabulary'],
  ['sound_detective', 'Sound Detective', 'Escucha y encuentra la información.', 'listening'],
  ['grammar_rescue', 'Grammar Rescue', 'Corrige una oración y explica tu idea.', 'grammar'],
  ['spell_master', 'Spell Master', 'Construye palabras letra a letra.', 'writing'],
  ['conversation_quest', 'Conversation Quest', 'Responde según la situación.', 'reading'],
  ['story_explorer', 'Story Explorer', 'Explora historias y comprende detalles.', 'reading'],
  ['word_rain', 'Word Rain', 'Atrapa palabras de la categoría correcta.', 'vocabulary'],
  ['english_rush', 'English Rush', 'Combina habilidades en un reto rápido.', 'mixed'],
]

export function answerType(kind, question = {}) {
  if (kind === 'english_rush') return question.response || 'choice'
  return ({ sentence_builder: 'order', memory_match: 'pairs', grammar_rescue: 'text', spell_master: 'letters', word_rain: 'multi' })[kind] || 'choice'
}

export function newQuestion(kind) {
  const id = () => `o-${crypto.randomUUID().replaceAll('-', '').slice(0, 24)}`
  const mode = answerType(kind)
  const labels = mode === 'pairs' ? ['English word', 'Significado', 'Another word', 'Otro significado'] : mode === 'order' ? ['I', 'like', 'apples.'] : mode === 'letters' ? ['a', 'p', 'p', 'l', 'e'] : mode === 'text' ? [] : ['Option one', 'Option two']
  const options = labels.map(label => ({ id: id(), label, image: '' }))
  return { prompt: '', explanation: '', passage: '', audio_key: '', response: kind === 'english_rush' ? 'choice' : '', options, solution: mode === 'pairs' ? [[options[0].id, options[1].id], [options[2].id, options[3].id]] : mode === 'order' ? options.map(o => o.id) : mode === 'letters' ? 'apple' : mode === 'text' ? '' : mode === 'multi' ? [options[0].id] : options[0].id }
}
