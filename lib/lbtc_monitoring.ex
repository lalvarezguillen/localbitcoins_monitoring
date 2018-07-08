defmodule LbtcMonitoring do
  @moduledoc """
  Documentation for LbtcMonitoring.
  """

  @doc """
  Hello world.

  ## Examples

      iex> LbtcMonitoring.hello
      :world

  """
  def startSearch(currency, keywords) do
    url = "https://localbitcoins.com/sell-bitcoins-online/#{currency}/.json"
    lower_kws = Enum.map(keywords, fn k -> String.downcase(k) end)
    getOffers(url, [], currency, lower_kws)
  end

  def getOffers(url, acc, currency, keywords) do
    IO.puts(url)
    resp = HTTPotion.get(url)

    if resp.status_code == 200 do
      [partialOffers, next] = parseResponse(resp.body)
      acc = acc ++ partialOffers

      if next do
        getOffers(next, acc, currency, keywords)
      else
        Enum.filter(acc, fn o -> checkIfInteresting(o, keywords) end)
      end
    end
  end

  def parseResponse(respBody) do
    resp = Jason.decode!(respBody)
    [resp["data"]["ad_list"], resp["pagination"]["next"]]
  end

  def checkIfInteresting(offer, keywords) do
    Enum.any?(keywords, fn k ->
      String.downcase(offer["data"]["msg"]) =~ k or
        String.downcase(offer["data"]["bank_name"]) =~ k
    end)
  end
end

defmodule LbtcMonitoring.CLI do
  def main(args) do
    cleanArgs = Enum.map(args, fn elem -> cleanArg(elem) end)
    IO.puts(cleanArgs)

    parsedArgs =
      OptionParser.parse(
        cleanArgs,
        strict: [currency: :string]
      )

    {[currency: curr], keywords, _} = parsedArgs
    IO.puts(curr)
    Enum.map(keywords, fn elem -> IO.puts(elem) end)
    searchLbtc(curr, keywords)
  end

  defp cleanArg(arg) do
    String.trim(arg)
    |> String.downcase()
  end

  defp searchLbtc(currency, keywords) do
    LbtcMonitoring.startSearch(currency, keywords)
    |> IO.puts()
  end
end
