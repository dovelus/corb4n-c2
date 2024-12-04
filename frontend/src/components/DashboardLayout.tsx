import * as React from "react";
import { Link, useNavigate, Outlet } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { Trash2, ExternalLink, RefreshCw } from "lucide-react";
import { Button } from "./ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "./ui/table";
import { Topbar } from "./Topbar";
import { fetchBeacons, Beacon } from "../api/beacons";

function DashboardLayoutComponent() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = React.useState(false);

  return (
    <div className="flex flex-col min-h-screen bg-background text-foreground">
      <Topbar
        onMobileMenuToggle={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
      />
      {isMobileMenuOpen && (
        <nav className="md:hidden bg-background border-b p-4">
          <div className="flex flex-col space-y-2">
            <Button asChild variant="ghost" size="sm">
              <Link to="/">Dashboard</Link>
            </Button>
            <Button asChild variant="ghost" size="sm">
              <Link to="/settings">Settings</Link>
            </Button>
            <Button asChild variant="ghost" size="sm">
              <Link to="/status">Status</Link>
            </Button>
          </div>
        </nav>
      )}
      <main className="flex-1 w-full py-6">
        <div className="container mx-auto px-4">
          <Outlet />
        </div>
      </main>
    </div>
  );
}

function Dashboard() {
  const navigate = useNavigate();
  const {
    data: beacons,
    isLoading,
    error,
    refetch,
    isFetching,
  } = useQuery<Beacon[], Error>({
    queryKey: ["beacons"],
    queryFn: fetchBeacons,
  });

  const deleteBeacon = (name: string) => {
    console.log(`Deleting beacon with id: ${name}`);
  };

  const viewBeaconDetails = (name: string) => {
    navigate(`/beacon/${name}`);
  };

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>An error occurred: {error.message}</div>;

  return (
    <>
      <h1 className="text-3xl font-bold mb-6 text-center">Beacon Dashboard</h1>
      <div className="overflow-x-auto">
        <div className="inline-block min-w-full align-middle">
          <div className="overflow-hidden border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="px-3 py-2">ID</TableHead>
                  <TableHead className="px-3 py-2">Beacon Name</TableHead>
                  <TableHead className="px-3 py-2">IP Address</TableHead>
                  <TableHead className="px-3 py-2">Status</TableHead>
                  <TableHead className="px-3 py-2 w-80">
                    <div className="flex items-center justify-between">
                      <span className="flex-grow text-center">Actions</span>
                      <Button
                        onClick={() => refetch()}
                        disabled={isFetching}
                        size="icon"
                        variant="ghost"
                        className="h-6 w-6 ml-2"
                      >
                        <RefreshCw
                          className={`h-4 w-4 ${isFetching ? "animate-spin" : ""}`}
                        />
                        <span className="sr-only">
                          {isFetching
                            ? "Refreshing beacons"
                            : "Refresh beacons"}
                        </span>
                      </Button>
                    </div>
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {beacons && beacons.length > 0 ? (
                  beacons.map((beacon, index) => (
                    <TableRow key={beacon.beaconID} className="h-10">
                      <TableCell className="font-medium px-3 py-2">
                        {index + 1}
                      </TableCell>
                      <TableCell className="px-3 py-2">
                        {beacon.beaconID}
                      </TableCell>
                      <TableCell className="px-3 py-2">
                        {beacon.beaconExtIP}
                      </TableCell>
                      <TableCell className="px-3 py-2">
                        <span
                          className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
                            calculateStatus(beacon.lastUpdate) === "online"
                              ? "bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100"
                              : "bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100"
                          }`}
                        >
                          {calculateStatus(beacon.lastUpdate)}
                        </span>
                      </TableCell>
                      <TableCell className="px-3 py-2">
                        <div className="flex justify-center space-x-2">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => viewBeaconDetails(beacon.beaconID)}
                            className="h-8 px-2 text-xs"
                          >
                            <ExternalLink className="h-3 w-3 mr-1" />
                            View Details
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => deleteBeacon(beacon.beaconID)}
                            className="h-8 px-2 text-xs"
                          >
                            <Trash2 className="h-3 w-3 mr-1" />
                            Delete
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                ) : (
                  <TableRow className="hover:bg-transparent">
                    <TableCell
                      colSpan={5}
                      className="h-24 text-center text-muted-foreground"
                    >
                      No Beacons Found
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </div>
      </div>
    </>
  );
}

function Settings() {
  return <h1 className="text-3xl font-bold mb-6 text-center">Settings Page</h1>;
}

function Status() {
  return <h1 className="text-3xl font-bold mb-6 text-center">Status Page</h1>;
}

function calculateStatus(lastUpdate: string) {
  const lastUpdateTime = new Date(lastUpdate);
  const currentTime = new Date();
  const timeDifference =
    (currentTime.getTime() - lastUpdateTime.getTime()) / (1000 * 60 * 60); // in hours
  return timeDifference < 6 ? "online" : "offline";
}

DashboardLayoutComponent.Dashboard = Dashboard;
DashboardLayoutComponent.Settings = Settings;
DashboardLayoutComponent.Status = Status;

export default DashboardLayoutComponent;
