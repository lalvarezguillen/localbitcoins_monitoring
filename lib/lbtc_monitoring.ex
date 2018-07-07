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

  @currency "VEF"
  @keywords ["bvc", "bvdc", "vdc", "venezolano"]

  def startSearch() do
    url = "https://localbitcoins.com/sell-bitcoins-online/#{@currency}/.json"
    getOffers(url, [])
  end

  def getOffers(url, acc) do
    IO.puts(url)
    resp = HTTPotion.get(url)
    [partialOffers, next] = parseResponse(resp.body)
    acc = acc ++ partialOffers

    if next do
      getOffers(next, acc)
    else
      Enum.filter(acc, fn o -> checkIfInteresting(o) end)
    end
  end

  def parseResponse(respBody) do
    resp = Jason.decode!(respBody)
    [resp["data"]["ad_list"], resp["pagination"]["next"]]
  end

  def checkIfInteresting(offer) do
    Enum.any?(@keywords, fn k -> offer["data"]["msg"] =~ k or offer["data"]["bank_name"] =~ k end)
  end
end
