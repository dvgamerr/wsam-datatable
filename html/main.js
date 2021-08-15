export default async ({ add }) => {
  const addResult = add(24, 24);
  document.body.textContent = `Hello World! addResult: ${addResult}`;
}
