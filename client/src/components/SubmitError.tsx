
interface SubmitErrorProps {
  err: string | null;
}


export default function SubmitError({err}:  SubmitErrorProps) {
    


    return (
        <div style={{ color: 'red' }}>{err}</div>
    )
} 