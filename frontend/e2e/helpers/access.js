import { expect } from '@playwright/test'

// Explicit expected answers for the test; never read live solutions from an API or DB.
const expected = {
  B2: new Set([
    'If I had known, I would have called.',
    'The association is worth investigating, but causation is not established.',
    'although',
    'A Monday delivery conditional on timely confirmation',
    'The manager said that they were checking the figures.',
    'The review praises the sound but questions the contextual guidance.',
    'The samples were tested yesterday.',
    'Bus passengers benefited; wider effects remain unmeasured.',
    'Could you specify the deadline and which task takes priority?',
    'might',
    'The headline generalises beyond the sample.',
    'Some prefer online study, whereas others value classroom discussion.',
    'Qualified support with a proposed trial',
    'They must have left.',
    'If you can provide forty now, we can accept the balance later.',
    'Ten percent is a target, not yet a confirmed annual result.',
  ]),
  C1: new Set([
    'Rarely have we encountered such a contradiction.',
    'The improvement is observed, but its attribution to flexible work alone remains uncertain.',
    'I appreciate the deadline, but I cannot approve publication until the figures are verified.',
    'The studies examine different groups and outcomes.',
    'The absence of accountability is observable; a deceptive motive would need more evidence.',
    'with',
    'The sufficiency of access alone',
    'If we had checked, we would know now.',
    'The daily average may conceal brief peaks; those require more detailed measurements.',
    'A review is intended, but no publication date is promised.',
    'By no means is this a complete explanation.',
    'The distinction between entry and participation',
    'The prediction remains conditional on resources the budget does not provide.',
    'We are not yet in a position to agree.',
    'It defeats the universal guarantee, not every possible benefit of openness.',
    'Test whether the decision remains workable across the range.',
  ]),
  C2: new Set([
    'Failure is already established and awaits admission.',
    'Let us specify the criteria and distinguish missed, met and unmeasured outcomes.',
    'As irony supported by the mismatch with events',
    'Universal replacement is excluded; the covered subset still needs clarification.',
    'A challenge has been presented, not a demonstrated disproof.',
    'A suggestion of reluctant yielding not established by the interview',
    'Their observations may coexist, while a genuine disagreement about significance remains.',
    'Whether the reforms are credible and will be implemented',
    'It may work with sufficient capacity; evidence in fragile settings remains limited.',
    'What we need is an explicit criterion.',
    'A rhetorical label substituting for a substantive justification',
    'Certainty is limited by the evidence, not necessarily by the importance of the question.',
    'Partial non-approval is not wholesale rejection of the report.',
    'as',
    'The mild wording downplays the established extent of the delay.',
    "Remove incidental procedure, but retain the conclusion's scope and uncertainty.",
  ]),
}

export async function unlockByExam(page, level) {
  await page.goto(`/mapa?level=${level}`)
  await expect(page.getByRole('heading', { name: `${level} todavía está bloqueado` })).toBeVisible()
  await page.getByRole('link', { name: `Examen de acceso a ${level}`, exact: true }).click()
  await page.getByRole('button', { name: 'Comenzar o retomar examen' }).click()
  await expect(page.locator('.question-card')).toHaveCount(12)
  const id = new URL(page.url()).pathname.split('/').pop()
  const before = (await (await page.request.get(`http://127.0.0.1:8080/api/v1/access-exams/${id}`)).json()).data
  expect(JSON.stringify(before)).not.toContain('"solution"')
  expect(JSON.stringify(before)).not.toContain('"explanation"')
  await page.reload()
  await expect(page.locator('.question-card')).toHaveCount(12)
  const after = (await (await page.request.get(`http://127.0.0.1:8080/api/v1/access-exams/${id}`)).json()).data
  expect(after.expires_at).toBe(before.expires_at)
  for (const question of await page.locator('.question-card').all()) {
    const labels = await question.getByRole('button').allTextContents()
    const correct = labels.filter(label => expected[level].has(label))
    expect(correct).toHaveLength(1)
    await question.getByRole('button', { name: correct[0], exact: true }).click()
  }
  await page.getByRole('button', { name: 'Enviar examen', exact: true }).click()
  await expect(page.getByRole('heading', { name: `Nivel ${level} desbloqueado` })).toBeVisible()
  await expect(page.locator('.result-score')).toHaveText('100/100')
  await page.reload()
  await expect(page.getByRole('heading', { name: `Nivel ${level} desbloqueado` })).toBeVisible()
}
