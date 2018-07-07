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
