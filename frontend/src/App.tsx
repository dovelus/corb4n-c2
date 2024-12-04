import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { ThemeProvider } from "./components/theme-provider";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import DashboardLayoutComponent from "./components/DashboardLayout";
import BeaconDetailsPage from "./pages/BeaconDetails";
import { Toaster } from "./components/ui/toaster";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
        <Router>
          <Routes>
            <Route path="/" element={<DashboardLayoutComponent />}>
              <Route index element={<DashboardLayoutComponent.Dashboard />} />
              <Route path="beacon/:id" element={<BeaconDetailsPage />} />
              <Route
                path="settings"
                element={<DashboardLayoutComponent.Settings />}
              />
              <Route
                path="status"
                element={<DashboardLayoutComponent.Status />}
              />
            </Route>
          </Routes>
        </Router>
        <Toaster />
      </ThemeProvider>
    </QueryClientProvider>
  );
}

export default App;
