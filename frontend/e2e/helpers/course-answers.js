import { readFileSync } from 'node:fs'

// Test-only authoring fixtures. Never shipped in the application and never read
// from a privileged API or live database. Exercises still use real start/submit.
const read = path => readFileSync(new URL(`../../../backend/${path}`, import.meta.url), 'utf8')
const course = JSON.parse(read('curriculum/course.json'))
const embedded = (path, variable) => {
  const match = read(path).match(new RegExp(`${variable} jsonb := '([\\s\\S]*?)'::jsonb`))
  if (!match) throw new Error(`Missing authoring fixture ${variable}`)
  return JSON.parse(match[1].replaceAll("''", "'"))
}
const starter = embedded('migrations/sql/000003_starter_content.up.sql', 'topics')
const extended = embedded('migrations/sql/000006_extended_practice.up.sql', 'packs')

export function expectedSubmission(lesson, activity) {
  const unit = course.find(u => u.level === lesson.level_code && u.title === lesson.unit_title)
  const old = starter.find(u => u.level === lesson.level_code && u.title === lesson.unit_title)
  const pack = extended.find(u => u.level === lesson.level_code && u.unit === lesson.unit_title)
  const game = pack?.games.find(g => g.kind === activity.kind)
  const answers = {}
  for (const question of activity.questions) {
    const option = label => {
      const found = question.options.find(o => o.label === label)
      if (!found) throw new Error(`Missing option ${label} in ${lesson.unit_title}`)
      return found.id
    }
    if (unit) {
      if (activity.kind === 'grammar_rescue') answers[question.id] = unit.grammar.find(q => q[0] === question.prompt)[1]
      else {
        const list = question.audio_key ? unit.listening_questions : [...unit.comprehension, ...unit.interaction]
        const expected = list.find(q => q[0] === question.prompt)
        if (!expected) throw new Error(`Missing course fixture: ${question.prompt}`)
        answers[question.id] = option(expected[1][expected[2]])
      }
    } else if (old) {
      if (activity.kind === 'word_hunter') answers[question.id] = option(old.words.find(w => `Find: ${w.en}` === question.prompt).es)
      else if (activity.kind === 'sentence_builder') answers[question.id] = old.sentences.find(s => s.es === question.prompt).words.map(option)
      else answers[question.id] = old.words.map(w => [option(w.en), option(w.es)])
    } else if (game) {
      const expected = game.questions.find(q => q.prompt === question.prompt && (!q.passage || q.passage === question.passage))
      if (!expected) throw new Error(`Missing extended fixture: ${question.prompt}`)
      const type = question.response || ({ grammar_rescue: 'text', spell_master: 'letters', word_rain: 'multi' })[activity.kind] || 'choice'
      if (type === 'text') answers[question.id] = expected.answer
      else if (type === 'order') answers[question.id] = expected.options.map(option)
      else if (type === 'multi') answers[question.id] = expected.answer.map(i => option(expected.options[i]))
      else if (type === 'letters') {
        const used = new Set()
        answers[question.id] = [...expected.answer].map(letter => {
          const tile = question.options.find(o => o.label === letter && !used.has(o.id))
          if (!tile) throw new Error('Missing letter tile')
          used.add(tile.id); return tile.id
        })
      } else answers[question.id] = option(expected.options[expected.answer])
    } else throw new Error(`No expected authoring fixture for ${lesson.unit_title}`)
  }
  return { answers }
}
