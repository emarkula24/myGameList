
interface SubmitErrorProps {
  err: string | null;
}


export default function SubmitError({ err}: SubmitErrorProps) {

  console.log(err)
  return (
    <div style={{ textAlign: "center", color: 'red' , fontSize: "2em", width: "20%", height: "10%"}}>{err}</div>
  )
} 