import React, { useMemo } from 'react';
import { useAppSelector } from '../../app/hooks';
import { selectAllRequests } from './logsSlice';
import FilterForm from './FilterForm';
import RequestTable from './LogTable';

const Logs = () => {
	const requestLogs = useAppSelector(selectAllRequests);
	const ipAddresses = useAppSelector(state => state.logs.ipAddresses);

	const filteredRequests = useMemo(() => {
		const selectedIps = ipAddresses
			.filter((ip) => ip.selected)
			.map((ip) => ip.value);

		return requestLogs.filter((request) => selectedIps.includes(request.ipAddress));
	}, [requestLogs, ipAddresses]);

	return (
		<>
			<FilterForm />
			{filteredRequests.length > 0
				? <RequestTable requests={filteredRequests} />
				: <div className='message'>No logs to render...</div>}
		</>
	);
};

export default Logs;