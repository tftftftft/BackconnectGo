<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from '$app/navigation';
  import { ConicGradient } from '@skeletonlabs/skeleton';
import type { ConicStop } from '@skeletonlabs/skeleton';

  interface ProxyInfo {
    ServerIP: string;
    ServerListeningPort: string;
    ProxyIP: string;
    CountryCode: string;
    Region: string;
    City: string;
    Zip: string;
    Mobile: boolean;
    Proxy: boolean;
    Hosting: boolean;
  }

  let proxies: ProxyInfo[] = [];
  let selectedCountry = "";
  let searchRegion = "";
  let searchCity = "";
  let searchZip = "";
  let searchIP = "";

  const fetchProxies = async () => {
    try {
      const res = await fetch('http://127.0.0.1:30000/fetchProxies', {
      credentials: 'include'
      });
      if (!res.ok) {
        throw new Error("Network response was not ok");
      }
      const result = await res.json();
      if (result.data && Array.isArray(result.data)) {
        proxies = result.data as ProxyInfo[];
      }
    } catch (error) {
      console.error("There was a problem with the fetch operation:", error);
    }
  };


  onMount(() => {
  // Check for session data
  const session = localStorage.getItem('session');
  if (!session) {
    goto('/');
    return;
  }
  
  
  
  // Fetch proxy data
  fetchProxies();
  const interval = setInterval(fetchProxies, 30000);
  
  return () => {
    clearInterval(interval);
  };
});


$: filteredProxies = proxies.filter(proxy => 
  (!selectedCountry || proxy.CountryCode === selectedCountry) &&
  (!searchRegion || proxy.Region.toLowerCase().includes(searchRegion.toLowerCase())) &&
  (!searchCity || proxy.City.toLowerCase().includes(searchCity.toLowerCase())) &&
  (!searchZip || proxy.Zip.includes(searchZip)) &&
  (!searchIP || proxy.ProxyIP.includes(searchIP)) // Added this line for IP filtering
);

  $: availableCountryCodes = Array.from(new Set(proxies.map(proxy => proxy.CountryCode)));
  // Computed statistics
  $: totalProxies = proxies.length;
</script>


<!-- Statistics Panel -->
<div class="flex justify-between items-center p-4  rounded-lg mb-4  text-lg rounded-md bg-surface-700">
  <div class="stat-item">
    <span class="font-bold">Online Proxies:</span> {totalProxies}
  </div>

  <!-- Add more statistics as needed -->
</div>
<!-- Search Panel -->
<div class="flex justify-between items-center p-4 rounded-lg shadow-lg space-x-4 rounded-md bg-surface-700 mb-4">
  <!-- Dropdown for Country Codes -->
  <select bind:value={selectedCountry} class="select p-2 rounded-lg  focus:ring focus:ring-opacity-50">
    <option value="">Select Country</option>
    {#each availableCountryCodes as code}
      <option value={code}>{code}</option>
    {/each}
  </select>
     <!-- Text Input for IP Address -->
  <input bind:value={searchIP} class="input p-2 rounded-lg  focus:ring focus:ring-opacity-50" title="Input (text)" type="text" placeholder="Proxy IP"/>
  <!-- Text Input for Region -->
  <input bind:value={searchRegion} class="input p-2 rounded-lg  focus:ring focus:ring-opacity-50" title="Input (text)" type="text" placeholder="Region"/>
  <!-- Text Input for City -->
  <input bind:value={searchCity} class="input p-2 rounded-lg  focus:ring focus:ring-opacity-50" title="Input (text)" type="text" placeholder="City"/>
  <!-- Text Input for Zip -->
  <input bind:value={searchZip} class="input p-2 rounded-lg  focus:ring focus:ring-opacity-50" title="Input (text)" type="text" placeholder="Zip"/>
 

</div>
<!-- Table -->
<div class="table-container mb-8">
  <table class="table table-hover w-full">
    <thead>
      <tr>
        <th class="text-center">Server IP:Port</th> <!-- Modified column header -->
        <th class="text-center">Proxy IP</th>
        <th class="text-center">Country Code</th>
        <th class="text-center">Region</th>
        <th class="text-center">City</th>
        <th class="text-center">Zip</th>
      </tr>
    </thead>
    <tbody>
      {#if Array.isArray(filteredProxies)}
        {#each filteredProxies as proxy (proxy.ServerIP)}
          <tr>
            <td class="text-center">{proxy.ServerIP}:{proxy.ServerListeningPort}</td> <!-- Combined fields -->
            <td class="text-center">{proxy.ProxyIP}</td>
            <td class="text-center">{proxy.CountryCode}</td>
            <td class="text-center">{proxy.Region}</td>
            <td class="text-center">{proxy.City}</td>
            <td class="text-center">{proxy.Zip}</td>
          </tr>
        {/each}
      {/if}
    </tbody>
  </table>
</div>