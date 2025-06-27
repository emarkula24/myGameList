
export function Spinner({
  show,
  wait,
}: {
  show?: boolean
  wait?: `delay-${number}`
}) {
    const isVisible = show ?? true
  return (
    <div 
    className={`spinner ${isVisible ? 'visible' : 'hidden'} ${wait ?? 'delay-300'}`}
    >
      ⍥
    </div>
  )
}