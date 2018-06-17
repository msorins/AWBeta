import { Status } from "./Status";

export interface ISolver {
    getStatuses: () => Status[],
    checkIfStatusHasChanged: () => Boolean
    solve: () => void
}