import { AuthProvider } from "@ssoeasy-dev/react";
import { AppRouter } from "./router";

function App() {
  return (
    <AuthProvider
      config={{
        redirectUri: "http://localhost:5174/",
        loginPath: "login",
        serviceId: "97015f75-5899-4143-b12a-4becfbadb56d",
      }}
    >
      <AppRouter />
    </AuthProvider>
  );
}

export default App;
