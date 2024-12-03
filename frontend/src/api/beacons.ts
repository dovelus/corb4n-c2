import axios from "axios";

export interface Beacon {
  beaconID: string;
  hostname: string;
  beaconIntIP: string;
  beaconExtIP: string;
  os: string;
  endpointProt: string;
  processID: number;
  processUser: string;
  lastUpdate: string;
}

export const fetchBeacons = async (): Promise<Beacon[]> => {
  const response = await axios.get<Beacon[]>("/api/beacons");
  return response.data;
};

export const fetchBeaconById = async (id: string): Promise<Beacon> => {
  const response = await axios.get<Beacon[]>("/api/beacons");
  const beacon = response.data.find((b) => b.beaconID === id);
  if (!beacon) {
    throw new Error(`Beacon with id ${id} not found`);
  }
  return beacon;
};
