import '../styles/globals.css'
import type { AppProps } from 'next/app'
import { ApolloClient, ApolloProvider, HttpLink, InMemoryCache } from "@apollo/client";
import { FindTaskDocument, TasksDocument, useFindTaskQuery } from '../src/generated/graphql';

export const createApolloClient = () => {
  return new ApolloClient({
    link: new HttpLink({
      uri: 'http://localhost:3000/query',
    }),
    cache: new InMemoryCache(),
  });
 };
const client = createApolloClient()

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <ApolloProvider client={client}>
      <Component {...pageProps} />
    </ApolloProvider>
  )
}
export default MyApp
