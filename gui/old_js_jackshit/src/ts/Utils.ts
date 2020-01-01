export function conditionalStr(cnd: string | number | boolean | object, value: string, replacement?: string): string {
    const rep = replacement ? replacement : ''

    if (typeof cnd == 'string') {
        return cnd !== '' ? value : rep
    } else if (typeof cnd == 'number') {
        return cnd !== 0 ? value : rep
    } else if (typeof cnd == 'boolean') {
        return cnd ? value : rep
    } else {
        return cnd != null && cnd != undefined ? value : rep
    }
}
