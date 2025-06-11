import { RegisterForm } from "@/components/registerForm"

export default function Page() {
  return (
    <div className="flex w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <RegisterForm userType="seller"  />
      </div>
    </div>
  )
}